package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/medvedevse/kafka-hotels-book/internal/entity"
	"github.com/medvedevse/kafka-hotels-book/internal/repository/redis"
)

type Handler struct {
	idempotencyService *redis.Service
}

func NewHandler(redisAddr string) (*Handler, error) {
	service, err := redis.NewService(redisAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create idempotency service: %w", err)
	}

	return &Handler{idempotencyService: service}, nil
}

// HandlerMessage обрабатывает сообщение с проверкой идемпотентности
func (h *Handler) HandlerMessage(msg []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Kafka message: %s", string(msg))

	// Парсим сообщение для получения уникального id
	messageKey, err := h.extractMessageKey(msg)
	if err != nil {
		log.Printf("Failed to extract message key: %v", err)
		return err
	}

	// Проверяем, было ли сообщение уже обработано
	processed, err := h.idempotencyService.IsProcessed(ctx, messageKey)
	if err != nil {
		log.Printf("Failed to check if message was processed: %v", err)
		return err
	}

	if processed {
		log.Printf("Message with key %s already processed, skipping", messageKey)
		return nil
	}

	// Отмечаем сообщение как обработанное
	if err := h.idempotencyService.MarkAsProcessed(ctx, messageKey, string(msg)); err != nil {
		log.Printf("Failed to mark message as processed: %v", err)
		// Не возвращаем ошибку, так как сообщение уже обработано
	}
	return nil
}

// extractMessageKey извлекает уникальный ключ из сообщения
func (h *Handler) extractMessageKey(msg []byte) (string, error) {
	var hotelBook entity.HotelBook
	if err := json.Unmarshal(msg, &hotelBook); err != nil {
		// Если не удается распарсить как HotelBook, используем хеш самого сообщения
		return h.hashMessage(msg), nil
	}

	// Используем ID как ключ идемпотентности
	if hotelBook.Id != "" {
		return hotelBook.Id, nil
	}

	// Если ID нет, создаем ключ из сообщения
	return h.hashMessage(msg), nil
}

// Хеш для использования в качестве ключа
func (h *Handler) hashMessage(msg []byte) string {
	hash := sha256.Sum256(msg)
	return hex.EncodeToString(hash[:])
}

func (h *Handler) Close() error {
	if h.idempotencyService != nil {
		return h.idempotencyService.Close()
	}
	return nil
}
