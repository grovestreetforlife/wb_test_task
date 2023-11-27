package application

import (
	"context"
	"github.com/dany-ykl/tracer"
	"github.com/pkg/errors"
	"wb_test_task/api/internal/config"
	"wb_test_task/api/internal/delivery/http"
	"wb_test_task/api/internal/services"
	psql "wb_test_task/api/internal/storage/psql"
	"wb_test_task/api/internal/storage/redis"
)

type Application struct {
	cfg          *config.Config
	httpServer   *http.Server
	postgres     *psql.Storage
	cache        *redis.Cache
	service      *services.Service
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
		return &Application{}, errors.Wrap(err, "fail to init postgres database")
	}

	service := services.New(services.Depends{
		OrderStorage: postgres.OrderStorage,
		OrderCache:   cache.OrderCache,
	})

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

	return &Application{
		cfg:          cfg,
		httpServer:   http.New(cfg.Server.HttpServer, service),
		postgres:     postgres,
		cancelTracer: cancelTracer,
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	if err := a.httpServer.Start(); err != nil {
		return errors.Wrap(err, "fail to start http server")
	}

	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "fail to shutdown http server")
	}
	return nil
}
