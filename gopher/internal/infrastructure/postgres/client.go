package postgres

import (
	"context"
	"gopher/internal/config"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDatabese(config config.Config) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), config.DatabaseURL())
	if err != nil {
		log.Fatalf("postgres dont connect: %s", err.Error())
	}
	return pool
}
