package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/medvedevse/kafka-hotels-book/config"
	"github.com/medvedevse/kafka-hotels-book/internal/entity"
	"github.com/medvedevse/kafka-hotels-book/internal/kafka"
)

func main() {
	cfg := config.NewConfig()

	p, err := kafka.NewProducer(cfg.Addr)
	if err != nil {
		log.Fatal(err)
	}

	keys := kafka.GenerateKeys()
	keyIdx := 0

	for {
		keyIdx++
		hotelName := fmt.Sprintf("TestHotel-%d", keyIdx)
		dateEnd := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+rand.Intn(10), 0, 0, 0, 0, time.UTC)
		msg := entity.NewHotelBook(hotelName, dateEnd).ToString()

		if err := p.Produce(msg, cfg.Topic, keys[keyIdx%len(keys)]); err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second * 3)
	}
}
