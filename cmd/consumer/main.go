package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/medvedevse/kafka-hotels-book/config"
	"github.com/medvedevse/kafka-hotels-book/internal/handler"
	"github.com/medvedevse/kafka-hotels-book/internal/kafka"
)

func main() {
	cfg := config.NewConfig()
	h := handler.NewHandler()

	idemHandler, err := handler.NewIdempotentHandler(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to create idempotent handler: %v", err)
	}
	defer func() {
		if closeErr := idemHandler.Close(); closeErr != nil {
			log.Printf("Error closing handler: %v", closeErr)
		}
	}()

	c, err := kafka.NewConsumer(h, idemHandler, cfg.Addr, cfg.Topic, cfg.ConsumerGroup)
	if err != nil {
		log.Fatal(err)
	}

	go c.Start()

	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM)
	<-sChan
	log.Fatal(c.Stop())
}
