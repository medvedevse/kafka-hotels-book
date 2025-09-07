package kafka

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	// sessionTimeout  = 30_000
	consumerTimeout = -1
)

type Handler interface {
	HandlerMessage(message []byte) error
}

type IdemHandler interface {
	HandlerMessage(msg []byte) error
}

type Consumer struct {
	consumer    *kafka.Consumer
	handler     Handler
	idemHandler IdemHandler
	stop        bool
}

func NewConsumer(handler Handler, idemHandler IdemHandler, addr []string, topic, consumerGroup string) (*Consumer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(addr, ","),
		"group.id":          consumerGroup,
		//"session.timeout.ms": sessionTimeout,
		// автосохранение офсета последнего сообщения
		// в локальную память консьюмера
		"enable.auto.offset.store": false,
		"auto.commit.interval.ms":  5000,
		// сохранение офсета в кафке
		"enable.auto.commit": true,
		// будем читать сообщения с начала партиции
		"auto.offset.reset": "earliest",
	}

	c, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, fmt.Errorf("error with new consumer: %w", err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:    c,
		handler:     handler,
		idemHandler: idemHandler,
	}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}

		msg, err := c.consumer.ReadMessage(consumerTimeout)
		if err != nil {
			log.Println(err)
		}

		if msg == nil {
			continue
		}

		// обработка сообщений
		if err := c.handler.HandlerMessage(msg.Value); err != nil {
			log.Print(err)
			continue
		}

		if err := c.idemHandler.HandlerMessage(msg.Value); err != nil {
			log.Print(err)
			continue
		}

		// вручную сохраняем офсет локально
		if _, err := c.consumer.StoreMessage(msg); err != nil {
			log.Print(err)
			continue
		}

		time.Sleep(time.Second * 3)
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	return c.consumer.Close()
}
