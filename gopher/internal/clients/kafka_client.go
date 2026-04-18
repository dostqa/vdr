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

func (s *KafkaService) StartConsume(ctx context.Context, topic, groupID string, handler func(kafka.Message)) {
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
		handler(msg)
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
