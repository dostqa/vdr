package database

import (
	"context"
	"fmt"

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
