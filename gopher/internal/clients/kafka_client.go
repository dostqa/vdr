package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"sync"

	"github.com/segmentio/kafka-go"
)

type ObjectPdn struct {
	Text      string  `json:"text"`       // маппинг на "text"
	Type      string  `json:"type"`       // маппинг на "type"
	StartTime float64 `json:"start_time"` // маппинг на "start_time"
	EndTime   float64 `json:"end_time"`   // маппинг на "end_time"
}

// OutputMessage основная структура сообщения
type OutputMessage struct {
	RequestID    int         `json:"request_id"`    // маппинг на "request_id"
	OldFilePath  string      `json:"old_file_path"` // маппинг на "old_file_path"
	NewFilePath  string      `json:"new_file_path"` // маппинг на "new_file_path"
	OriginalText string      `json:"original_text"` // маппинг на "original_text"
	AnonText     string      `json:"anon_text"`     // маппинг на "anon_text"
	ObjectsPdns  []ObjectPdn `json:"objects_pdns"`  // маппинг на "objects_pdns"
}
type InputMessage struct {
	RequestID int    `json:"request_id"`
	FilePath  string `json:"file_path"`
}

type KafkaService struct {
	Writer  *kafka.Writer
	Brokers []string
	readers []*kafka.Reader
	mu      sync.Mutex
}

func NewKafkaService(brokers []string) *KafkaService {
	return &KafkaService{
		Brokers: brokers,
		Writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (s *KafkaService) SendJSON(ctx context.Context, topic string, key string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	slog.Info(string(payload), key)

	// Отправка через уже существующий метод или напрямую через Writer
	if err := s.Writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}); err != nil {
		return fmt.Errorf("отправка в топик: %w", err)
	}

	return nil
}

func (s *KafkaService) StartConsume(ctx context.Context, topic, groupID string, handler func(ctx context.Context, msg OutputMessage)) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: s.Brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	s.mu.Lock()
	s.readers = append(s.readers, reader)
	s.mu.Unlock()

	defer reader.Close()

	log.Printf("Started consuming topic: %s", topic)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("error reading from %s: %v", topic, err)
			continue
		}

		// 1. Создаем переменную типа нашей структуры
		var outputMsg OutputMessage

		// 2. Распаковываем JSON из тела сообщения Kafka (msg.Value)
		if err := json.Unmarshal(msg.Value, &outputMsg); err != nil {
			log.Printf("could not decode message: %v", err)
			// Важно: решаем, пропускать ли битое сообщение или падать
			continue
		}

		// 3. Передаем уже готовую структуру в обработчик
		handler(ctx, outputMsg)
	}
}

func (s *KafkaService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.Writer.Close(); err != nil {
		return err
	}

	for _, r := range s.readers {
		r.Close()
	}
	return nil
}
