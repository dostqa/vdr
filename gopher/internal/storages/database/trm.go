package database

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct {
	manager *manager.Manager
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		manager: manager.Must(trmpgx.NewDefaultFactory(pool)),
	}
}

func (tm *TransactionManager) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.manager.Do(ctx, fn)
}
