package base

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) GetDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		return tx
	}
	return r.DB.WithContext(ctx)
}

func (r *BaseRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var models []T
	result := r.GetDB(ctx).Find(&models)
	return models, result.Error
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var model T
	result := r.GetDB(ctx).First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &model, nil
}

func (r *BaseRepository[T]) Create(ctx context.Context, model *T) error {
	return r.GetDB(ctx).Create(model).Error
}

func (r *BaseRepository[T]) Update(ctx context.Context, model *T) error {
	return r.GetDB(ctx).Save(model).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var model T
	return r.GetDB(ctx).Delete(&model, id).Error
}
