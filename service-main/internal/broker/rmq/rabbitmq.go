package rmq

import (
	"context"
	"encoding/json"

	"github.com/Dsmit05/ot-example/service-main/internal/models"
	"github.com/Dsmit05/ot-example/utils"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type Producer struct {
	queueName string
	ch        *amqp.Channel
}

func NewProducer(uri, queueName string) (*Producer, error) {
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

	return &Producer{
		queueName: q.Name,
		ch:        ch,
	}, nil
}

func (p *Producer) SendMsg(ctx context.Context, msg models.Message) error {
	msgB, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "json.Marshal() error")
	}

	tr := otel.Tracer("publish")

	amqpContext, messageSpan := tr.Start(ctx, "amqp_publish")
	defer messageSpan.End()

	headers := utils.InjectAMQPHeaders(amqpContext)

	err = p.ch.PublishWithContext(ctx,
		"",          // exchange
		p.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			Headers:     headers,
			ContentType: "application/json",
			Body:        msgB,
		})
	if err != nil {
		return errors.Wrap(err, "p.ch.PublishWithContext() error")
	}

	return nil
}
