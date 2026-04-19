package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gopher/internal/infrastructure/kafka/message"
)

type RequestRepository struct {
	*Repository
}

func NewRequestRepository(repo *Repository) *RequestRepository {
	return &RequestRepository{repo}
}

func (r *RequestRepository) Save(ctx context.Context) (int, error) {
	const op = "requestRepository.Save"

	var id int

	err := r.db(ctx).QueryRow(ctx,
		`INSERT INTO requests DEFAULT VALUES RETURNING id`,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *RequestRepository) IsReady(ctx context.Context, requestID int) (bool, error) {
	const op = "requestRepository.IsReady"

	var isReady bool

	err := r.db(ctx).QueryRow(ctx,
		`SELECT status FROM requests WHERE id = $1`,
		requestID,
	).Scan(&isReady)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isReady, nil
}

func (r *RequestRepository) SaveOutputMessage(ctx context.Context, msg message.OutputMessage) error {
	const op = "requestRepository"

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := `
        UPDATE requests
        SET status = TRUE,
            payload = $1
        WHERE id = $2`

	_, err = r.db(ctx).Exec(ctx, query, jsonData, msg.RequestID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RequestRepository) OutputMessage(ctx context.Context, requestID int) (message.OutputMessage, error) {
	const op = "requestRepository.OutputMessage"

	var jsonData []byte

	query := "SELECT payload FROM requests WHERE id = $1"
	err := r.db(ctx).QueryRow(ctx, query, requestID).Scan(&jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			return message.OutputMessage{}, fmt.Errorf("%s: request %d not found", op, requestID)
		}
		return message.OutputMessage{}, fmt.Errorf("%s: %w", op, err)
	}

	if len(jsonData) == 0 {
		return message.OutputMessage{}, fmt.Errorf("%s: payload is empty for request %d", op, requestID)
	}

	var msg message.OutputMessage
	if err := json.Unmarshal(jsonData, &msg); err != nil {
		return message.OutputMessage{}, fmt.Errorf("%s: failed to unmarshal payload: %w", op, err)
	}

	return msg, nil
}
