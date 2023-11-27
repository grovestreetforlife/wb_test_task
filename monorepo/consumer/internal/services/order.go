package services

import (
	"context"
	"wb_test_task/consumer/internal/domain"
	"wb_test_task/libs/model"
)

//go:generate mockgen -source=order.go -destination=mocks/mock.go
type orderStorage interface {
	Create(ctx context.Context, request *domain.OrderCreateRequest) (*model.Order, error)
}

type orderCache interface {
	Set(ctx context.Context, key string, order *model.Order) error
}

type orderService struct {
	store orderStorage
	cache orderCache
}

func newOrderService(store orderStorage, cache orderCache) *orderService {
	return &orderService{store: store, cache: cache}
}

// Create создание заказа
func (o *orderService) Create(ctx context.Context, request *domain.OrderCreateRequest) (*model.Order, error) {
	order, err := o.store.Create(ctx, request)
	if err != nil {
		return &model.Order{Items: []*model.Product{}}, err
	}

	if err := o.cache.Set(ctx, request.OrderUid, order); err != nil {
		return &model.Order{Items: []*model.Product{}}, err
	}

	return order, nil
}
