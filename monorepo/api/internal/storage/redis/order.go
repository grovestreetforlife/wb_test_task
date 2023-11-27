package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dany-ykl/tracer"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"time"
	"wb_test_task/api/internal/common"
	"wb_test_task/api/internal/domain"
	"wb_test_task/libs/model"
)

type orderCache struct {
	conn      *redis.Client
	ttlSecond int
}

const orderObjectPrefix = "orders"

func newOrderCache(conn *redis.Client, ttlSecond int) *orderCache {
	return &orderCache{conn: conn, ttlSecond: ttlSecond}
}

// GetByID вернуть order по id
func (o *orderCache) GetByID(ctx context.Context, id string) (*model.Order, error) {
	ctx, span := tracer.StartTrace(ctx, "redis-cache-get-order-by-id")
	span.SetAttributes(attribute.String("order-id", id))
	defer span.End()

	cmd := o.conn.Get(ctx, fmt.Sprintf("%s:%s", orderObjectPrefix, id))

	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return &model.Order{Items: []*model.Product{}}, common.WrapError{Err: domain.ErrOrderNotExists, Msg: domain.ErrOrderNotExists.Error()}
		}

		return &model.Order{Items: []*model.Product{}}, common.WrapError{Err: err, Msg: "fail to get order by id from cache"}
	}

	data, err := cmd.Bytes()
	if err != nil {
		return &model.Order{Items: []*model.Product{}}, common.WrapError{Err: err, Msg: "fail to get bytes"}
	}

	var order model.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return &model.Order{Items: []*model.Product{}}, common.WrapError{Err: err, Msg: "fail to unmarshal order"}
	}

	return &order, nil
}

// Set вставить order в redis cache
func (o *orderCache) Set(ctx context.Context, key string, order *model.Order) error {
	ctx, span := tracer.StartTrace(ctx, "redis-cache-set-order")
	span.SetAttributes(attribute.String("key", key))
	defer span.End()

	data, err := json.Marshal(order)
	if err != nil {
		return common.WrapError{Err: err, Msg: "fail to unmarshal order"}
	}

	cmd := o.conn.Set(ctx, fmt.Sprintf("%s:%s", orderObjectPrefix, key), data, time.Duration(o.ttlSecond)*time.Second)
	if err := cmd.Err(); err != nil {
		return common.WrapError{Err: err, Msg: "fail to set order in cache"}
	}

	return nil
}
