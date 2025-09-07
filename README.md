## Реализация producer и consumer для очередей бронирования отелей
### Билд
`docker-compose build`<br>
`docker-compose up -d`<br>

### Запуск
`docker-compose down -v`<br>
`docker-compose up -d`<br>
`go run ./cmd/producer` - запуск продюсера<br>
`go run ./cmd/consumer` - запуск консьюмера<br>

### UI
- **Redis Commander** доступен на `http://localhost:8081` для просмотра ключей
- **Kafka UI** доступен на `http://localhost:9020`

### Идемпотентность
1. **Producer отправляет сообщения** с уникальными ID
2. **Consumer получает сообщение** и проверяет в Redis, было ли оно уже обработано
3. **Если сообщение новое** - обрабатывается и помечается в Redis
4. **Если сообщение дублированное** - пропускается

### Библиотеки для работы с kafka
`github.com/confluentinc/confluent-kafka-go`<br>
`github.com/IBM/sarama`<br>
`github.com/segmentio/kafka-go`<br>
`github.com/lovoo/goka`<br>
