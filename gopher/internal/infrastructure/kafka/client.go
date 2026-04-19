package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"gopher/internal/infrastructure/kafka/message"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Consume(ctx context.Context, msg message.OutputMessage)
}

type Kafka struct {
	Writer  *kafka.Writer
	Brokers []string
	readers []*kafka.Reader
	mu      sync.Mutex
}

func NewKafka(brokers []string) *Kafka {
	return &Kafka{
		Brokers: brokers,
		Writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (s *Kafka) Close() error {
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

func (s *Kafka) SendJSON(ctx context.Context, topic string, key string, data message.InputMessage) error {
	const op = "Kafka.SendJSON"

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.Writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Kafka) StartConsume(ctx context.Context, topic string, consumer Consumer) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: s.Brokers,
		Topic:   topic,
	})

	s.mu.Lock()
	s.readers = append(s.readers, reader)
	s.mu.Unlock()

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}

		var outputMsg message.OutputMessage

		if err := json.Unmarshal(msg.Value, &outputMsg); err != nil {
			continue
		}

		consumer.Consume(ctx, outputMsg)
	}
}
