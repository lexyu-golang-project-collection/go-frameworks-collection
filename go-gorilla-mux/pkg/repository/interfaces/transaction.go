package base

import (
	"context"

	"gorm.io/gorm"
)

type TxManager struct {
	DB *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{DB: db}
}

// WithTx 在事務中執行函數
func (tm *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}
