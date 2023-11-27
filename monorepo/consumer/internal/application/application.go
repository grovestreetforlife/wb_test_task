package application

import (
	"context"
	"github.com/dany-ykl/tracer"
	"github.com/pkg/errors"
	"wb_test_task/consumer/internal/config"
	"wb_test_task/consumer/internal/consumer"
	"wb_test_task/consumer/internal/services"
	"wb_test_task/consumer/internal/storage/psql"
	"wb_test_task/consumer/internal/storage/redis"
)

type Application struct {
	cfg          *config.Config
	psqlStore    *psql.Storage
	redisCache   *redis.Cache
	service      *services.Service
	natsConsumer *consumer.Consumer
	cancelTracer func(ctx context.Context)
}

func New(ctx context.Context, configFile string) (*Application, error) {
	cfg, err := config.New(configFile)
	if err != nil {
		return &Application{}, errors.Wrap(err, "fail to init config")
	}

	cache, err := redis.New(cfg.Cache.RedisCache)
	if err != nil {
		return &Application{}, errors.Wrap(err, "fail to init redis cache")
	}

	postgres, err := psql.New(ctx, cfg.Database.PostgresDatabase)
	if err != nil {
		return &Application{}, errors.Wrap(err, "fail to init psql storage")
	}

	// init tracer
	cancelTracer, err := tracer.New(&tracer.Config{
		ServiceName:              cfg.Jaeger.ServiceName,
		Host:                     cfg.Jaeger.Host,
		Port:                     cfg.Jaeger.Port,
		Environment:              cfg.Jaeger.Environment,
		TraceRatioFraction:       cfg.Jaeger.TraceRatioFraction,
		OTELExporterOTLPEndpoint: cfg.Jaeger.OTELExporterOTLPEndpoint,
	})
	if err != nil {
		return &Application{}, errors.Wrap(err, "fail to init tracer")
	}

	service := services.New(services.Depends{
		OrderStorage: postgres.OrderStorage,
		OrderCache:   cache.OrderCache,
	})

	natsConsumer, err := consumer.New(cfg.Consumer.NatsConsumer, service.OrderService)
	if err != nil {
		return &Application{}, errors.Wrap(err, "fail to init natsConsumer")
	}

	return &Application{
		cfg:          cfg,
		psqlStore:    postgres,
		service:      service,
		natsConsumer: natsConsumer,
		redisCache:   cache,
		cancelTracer: cancelTracer,
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	if err := a.natsConsumer.Start(ctx); err != nil {
		return errors.Wrap(err, "fail to start nats consumer")
	}
	return nil
}

func (a *Application) Shutdown(_ context.Context) error {
	if err := a.psqlStore.Shutdown(); err != nil {
		return err
	}
	return nil
}
