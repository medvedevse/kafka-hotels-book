package kafka

import (
	"errors"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

const (
	flushTimeout = 5000
	keysNumber   = 30
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(addr []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":                     strings.Join(addr, ","),
		"acks":                                  "all", // at least once
		"retries":                               5,
		"enable.idempotence":                    true, // Предотвращение дубликатов при retry
		"max.in.flight.requests.per.connection": 1,
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error with new producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message, topic, key string) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Key:   []byte(key),
	}

	deliveryChan := make(chan kafka.Event)

	if err := p.producer.Produce(msg, deliveryChan); err != nil {
		return err
	}

	e := <-deliveryChan
	switch err := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return err
	default:
		return errors.New("unknown event type")

	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}

func GenerateKeys() [keysNumber]string {
	var keys [keysNumber]string

	for i := 0; i < keysNumber; i++ {
		keys[i] = uuid.NewString()
	}

	return keys
}
