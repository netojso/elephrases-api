package gormx

import (
	"context"

	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	"gorm.io/gorm"
)

type txKey struct{}

// DB devolve a transação carregada no ctx, se houver, ou db com o ctx aplicado.
func DB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return db.WithContext(ctx)
}

type txManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) portrepository.TxManager {
	return &txManager{db: db}
}

// WithinTransaction abre uma transação (ou um savepoint, se já houver uma no ctx).
func (m *txManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return DB(ctx, m.db).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}
