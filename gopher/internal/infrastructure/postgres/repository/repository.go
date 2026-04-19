package repository

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func (r *Repository) db(ctx context.Context) trmpgx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.pool)
}
