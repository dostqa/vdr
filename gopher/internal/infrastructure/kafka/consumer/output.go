package consumer

import (
	"context"
	"fmt"
	"gopher/internal/infrastructure/kafka/message"
	"log/slog"
)

type RequestRepository interface {
	SaveOutputMessage(ctx context.Context, msg message.OutputMessage) error
}

type OutputConsumer struct {
	requestRepository RequestRepository
}

func NewOutputConsumer(repo RequestRepository) *OutputConsumer {
	return &OutputConsumer{
		requestRepository: repo,
	}
}

func (c *OutputConsumer) Consume(ctx context.Context, msg message.OutputMessage) {
	const op = "OutputConsumer.Consume"

	if err := c.requestRepository.SaveOutputMessage(ctx, msg); err != nil {
		slog.Error(fmt.Sprintf("%s: %s", op, err))
	}
}
