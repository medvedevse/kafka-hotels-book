package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Префикс для ключей
	idempotencyKeyPrefix = "idem:"
	// TTL для ключей
	defaultTTL = 24 * time.Hour
)

type Service struct {
	client *redis.Client
	ttl    time.Duration
}

func NewService(redisAddr string, ttl ...time.Duration) (*Service, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     os.Getenv("REDIS_PSWD"),
		DB:           0, // бд по умолчанию
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	serviceTTL := defaultTTL
	if len(ttl) > 0 {
		serviceTTL = ttl[0]
	}

	return &Service{
		client: client,
		ttl:    serviceTTL,
	}, nil
}

// IsProcessed проверяет, было ли сообщение уже обработано
func (s *Service) IsProcessed(ctx context.Context, messageKey string) (bool, error) {
	key := s.buildKey(messageKey)

	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check message existence: %w", err)
	}

	return exists > 0, nil
}

// MarkAsProcessed отмечает сообщение как обработанное
func (s *Service) MarkAsProcessed(ctx context.Context, messageKey, messageValue string) error {
	key := s.buildKey(messageKey)

	err := s.client.Set(ctx, key, messageValue, s.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to mark message as processed: %w", err)
	}

	return nil
}

// buildKey создает ключ для Redis из исходного ключа сообщения
func (s *Service) buildKey(messageKey string) string {
	hash := sha256.Sum256([]byte(messageKey))
	hashStr := hex.EncodeToString(hash[:])
	return idempotencyKeyPrefix + hashStr
}

func (s *Service) Close() error {
	return s.client.Close()
}
