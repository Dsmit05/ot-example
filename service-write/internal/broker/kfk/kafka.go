package kfk

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Dsmit05/ot-example/service-write/internal/models"
	"github.com/Dsmit05/ot-example/service-write/internal/repository"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
)

type MessageWithContext struct {
	context.Context
	models.Message
}

// Consumer represents a Sarama consumer group.
type Consumer struct {
	client  sarama.ConsumerGroup
	topic   string
	handler sarama.ConsumerGroupHandler
	r       repository.ReposytoryI
}

func NewConsumer(uri, topic string, r repository.ReposytoryI) (*Consumer, error) {
	config := sarama.NewConfig()

	client, err := sarama.NewConsumerGroup([]string{uri}, "message", config)
	if err != nil {
		return nil, errors.Wrap(err, "sarama.NewConsumerGroup() error")
	}

	consumer := Consumer{
		client: client,
		topic:  topic,
		r:      r,
	}

	handler := otelsarama.WrapConsumerGroupHandler(&consumer)
	consumer.handler = handler

	return &consumer, nil
}

func (c *Consumer) Process() {
	for {
		if err := c.client.Consume(context.Background(), []string{c.topic}, c.handler); err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg models.Message

		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Print(" Ошибка")
			continue
		}

		log.Printf("Message claimed: value = %+v, timestamp = %v, topic = %s", msg, message.Timestamp, message.Topic)

		// Todo: not working? missing parent span
		// Check https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/Shopify/sarama/otelsarama/example
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(message))

		tr := otel.Tracer("Consume")

		ctx, span := tr.Start(ctx, "Consume")

		if err := c.r.SetMsg(ctx, msg); err != nil {
			span.SetStatus(1, "db setMsg error")
		}

		session.MarkMessage(message, "")

		span.End()
	}

	return nil
}
