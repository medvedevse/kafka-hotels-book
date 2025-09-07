package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr          []string
	Topic         string
	ConsumerGroup string
	RedisAddr     string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	return &Config{
		Addr: []string{
			os.Getenv("KAFKA_ADDR_1"),
			os.Getenv("KAFKA_ADDR_2"),
			os.Getenv("KAFKA_ADDR_3"),
		},
		Topic:         os.Getenv("KAFKA_TOPIC"),
		ConsumerGroup: os.Getenv("KAFKA_CONSUMER_GROUP"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
	}
}
