package rmq

import (
	"context"
	"encoding/json"

	"github.com/Dsmit05/ot-example/service-write/internal/models"
	"github.com/Dsmit05/ot-example/service-write/internal/repository"
	"github.com/Dsmit05/ot-example/utils"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Consumer struct {
	msgs <-chan amqp.Delivery
	r    repository.ReposytoryI
}

func NewConsumer(uri, queueName string, r repository.ReposytoryI) (*Consumer, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, errors.Wrap(err, "amqp.Dial() error")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "conn.Channel() error")
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "ch.QueueDeclare() error")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, errors.Wrap(err, "ch.Consume() error")
	}

	return &Consumer{
		msgs: msgs,
		r:    r,
	}, nil
}

func (c *Consumer) Process() {
	for msgAMQP := range c.msgs {
		c.setInDB(msgAMQP)
	}
}

func (c *Consumer) setInDB(msgAMQP amqp.Delivery) error {
	tr := otel.Tracer("consume")
	ctx := utils.ExtractAMQPHeaders(context.Background(), msgAMQP.Headers)

	ctx, span := tr.Start(ctx, "consume msg")
	defer span.End()

	var msg models.Message
	if err := json.Unmarshal(msgAMQP.Body, &msg); err != nil {
		span.SetStatus(1, err.Error())
		return err
	}

	span.AddEvent("info", trace.WithAttributes(attribute.String("msg", msg.Msg)))

	if err := c.r.SetMsg(ctx, msg); err != nil {
		return errors.Wrap(err, "r.SetMsg error()")
	}

	return nil
}
