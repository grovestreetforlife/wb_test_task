package psql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"wb_test_task/api/internal/config"
)

type pool interface {
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type Storage struct {
	conn         *pgxpool.Pool
	OrderStorage *orderStorage
}

func New(ctx context.Context, cfg config.PostgresDatabase) (*Storage, error) {
	connectCfg, err := pgxpool.ParseConfig(cfg.Url)
	if err != nil {
		return nil, err
	}
	connectCfg.MaxConns = int32(cfg.MaxOpenConn)
	connectCfg.MaxConnIdleTime = time.Minute * time.Duration(cfg.MaxConnLife)

	conn, err := pgxpool.NewWithConfig(ctx, connectCfg)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &Storage{
		conn:         conn,
		OrderStorage: newOrderStorage(conn),
	}, nil
}

func (p *Storage) Shutdown() error {
	p.conn.Close()
	return nil
}
