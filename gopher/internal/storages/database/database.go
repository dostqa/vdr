package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gopher/internal/clients"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type config interface {
	DataBaseURL() string
}

func InitDataBasePool(config config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), config.DataBaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %s", err.Error())
	}
	return pool, nil
}

type DataBase struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewDataBase(pool *pgxpool.Pool) *DataBase {
	return &DataBase{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func NewDataBaseWithConfig(config config) (*DataBase, error) {
	pool, err := InitDataBasePool(config)
	if err != nil {
		return nil, err
	}
	return NewDataBase(pool), nil
}

func (db *DataBase) db(ctx context.Context) trmpgx.Tr {
	return db.getter.DefaultTrOrDB(ctx, db.pool)
}

func (db *DataBase) SaveRequest(ctx context.Context) (int64, error) {
	const op = "database.CreateRequest"

	var id int64

	err := db.db(ctx).QueryRow(ctx,
		`INSERT INTO requests DEFAULT VALUES RETURNING id`,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (db *DataBase) SaveFile(ctx context.Context, requestID int64, filepath string) (int64, error) {
	const op = "database.CreateFile"

	var id int64

	err := db.db(ctx).QueryRow(ctx,
		`INSERT INTO files (request_id, filepath)
		 VALUES ($1, $2)
		 RETURNING id`,
		requestID, filepath,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (db *DataBase) UpdateStatus(ctx context.Context, requestID int64) error {
	const op = "database.UpdateStatus"

	_, err := db.db(ctx).Exec(ctx, `UPDATE requests SET status=true WHERE id=$1`, requestID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (db *DataBase) IsRequestReady(ctx context.Context, requestID int64) (bool, error) {
	const op = "database.IsFileProcessed"

	var isReady bool

	err := db.db(ctx).QueryRow(ctx,
		`SELECT status FROM requests
	 	 WHERE id = $1`,
		requestID,
	).Scan(&isReady)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isReady, nil
}

func (db *DataBase) HandleRequestJsonb(ctx context.Context, msg clients.OutputMessage) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	query := `
        UPDATE requests
        SET status = TRUE,
            payload = $1
        WHERE id = $2`

	_, err = db.db(ctx).Exec(ctx, query, jsonData, msg.RequestID)
	if err != nil {
		return fmt.Errorf("failed to update request: %v", err)
	}

	return nil
}

func (db *DataBase) GetRequestJsonb(ctx context.Context, requestID int) (*clients.OutputMessage, error) {
	var jsonData []byte

	query := "SELECT payload FROM requests WHERE id = $1"
	err := db.db(ctx).QueryRow(ctx, query, requestID).Scan(&jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("request %d not found", requestID)
		}
		return nil, err
	}

	if len(jsonData) == 0 {
		return nil, fmt.Errorf("payload is empty for request %d", requestID)
	}

	var msg clients.OutputMessage
	if err := json.Unmarshal(jsonData, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	return &msg, nil
}
