package kfk

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/Dsmit05/ot-example/service-main/internal/models"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(uri, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	var err error

	producer, err := sarama.NewSyncProducer(strings.Split(uri, ","), config)
	if err != nil {
		return nil, errors.Wrap(err, "sarama.NewSyncProducer() error")
	}

	producer = otelsarama.WrapSyncProducer(config, producer)

	return &Producer{producer: producer, topic: topic}, nil
}

func (p *Producer) SendMsg(ctx context.Context, msg models.Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "json.Marshal error")
	}

	kMsg := sarama.ProducerMessage{Topic: p.topic, Key: sarama.StringEncoder("message"),
		Value: sarama.ByteEncoder(b)}

	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(&kMsg))

	partition, offset, err := p.producer.SendMessage(&kMsg)
	if err != nil {
		return errors.Wrap(err, "p.producer.SendMessage() error")
	}

	log.Printf("data is stored on partition: %v, with offset: %v\n", partition, offset)

	return nil
}

func (p *Producer) Close() {
	p.producer.Close()
}
