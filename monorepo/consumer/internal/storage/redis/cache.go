package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"wb_test_task/consumer/internal/config"
)

type Cache struct {
	conn       *redis.Client
	cfg        *config.RedisCache
	OrderCache *orderCache
}

func New(cfg config.RedisCache) (*Cache, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
	})

	if cmd := conn.Ping(context.Background()); cmd.Err() != nil {
		return &Cache{}, cmd.Err()
	}

	return &Cache{
		conn:       conn,
		cfg:        &cfg,
		OrderCache: newOrderCache(conn, cfg.TtlSecond),
	}, nil
}

func (c *Cache) Shutdown() error {
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}
