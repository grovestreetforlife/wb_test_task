package services

import (
	"context"
	"github.com/dany-ykl/logger"
	"github.com/dany-ykl/tracer"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"time"
	"wb_test_task/api/internal/domain"
	"wb_test_task/libs/model"
)

//go:generate mockgen -source=order.go -destination=mocks/mock.go
type orderStorage interface {
	GetByID(ctx context.Context, id string) (*model.Order, error)
}

type orderCache interface {
	GetByID(ctx context.Context, id string) (*model.Order, error)
	Set(ctx context.Context, key string, order *model.Order) error
}

type orderService struct {
	store orderStorage
	cache orderCache
}

func newOrderService(store orderStorage, cache orderCache) *orderService {
	return &orderService{
		store: store,
		cache: cache,
	}
}

// GetByID вернуть order по id
func (o *orderService) GetByID(ctx context.Context, id string) (*model.Order, error) {
	ctx, span := tracer.StartTrace(ctx, "service-get-order-dy-id")
	span.SetAttributes(attribute.String("order-id", id))
	defer span.End()

	order, err := o.cache.GetByID(ctx, id)
	if !errors.Is(err, domain.ErrOrderNotExists) && err != nil {
		logger.Warn("service: fail to get order from cache", zap.Error(err))
	}

	if len(order.OrderUid) == 0 {
		order, err = o.store.GetByID(ctx, id)
		if err != nil {
			return &model.Order{Items: []*model.Product{}}, err
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := o.cache.Set(ctx, order.OrderUid, order); err != nil {
				logger.Warn("service: fail to set order in redis cache", zap.Error(err))
			}
		}()
	}

	return order, nil
}
