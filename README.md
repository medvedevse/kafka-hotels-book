# Библиотеки для работы с kafka
`github.com/confluentinc/confluent-kafka-go`
`github.com/IBM/sarama`
`github.com/segmentio/kafka-go`
`github.com/lovoo/goka`

# Билд
`docker-compose build`
`docker-compose up -d`

# Запустить
`docker-compose down -v`
`docker-compose up -d`
`go run ./cmd/producer` - запуск продюсера
`go run ./cmd/consumer` - запуск консьюмера

# Как работает идемпотентность
1. **Producer отправляет сообщения** с уникальными ID
2. **Consumer получает сообщение** и проверяет в Redis, было ли оно уже обработано
3. **Если сообщение новое** - обрабатывается и помечается в Redis
4. **Если сообщение дублированное** - пропускается

# UI

- **Redis Commander** доступен на `http://localhost:8081` для просмотра ключей
- **Kafka UI** доступен на `http://localhost:9020`
