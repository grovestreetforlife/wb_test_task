package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"wb_test_task/consumer/internal/common"
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

// Set вставить order в redis cache
func (o *orderCache) Set(ctx context.Context, key string, order *model.Order) error {
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
