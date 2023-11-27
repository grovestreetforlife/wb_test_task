package http

import (
	"context"
	"github.com/dany-ykl/logger"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"wb_test_task/api/internal/config"
	"wb_test_task/api/internal/delivery/http/health"
	"wb_test_task/api/internal/delivery/http/middleware"
	"wb_test_task/api/internal/delivery/http/v1api"
	"wb_test_task/api/internal/services"
)

type Server struct {
	server *echo.Echo
	cfg    config.HttpServer
}

func New(cfg config.HttpServer, service *services.Service) *Server {
	server := echo.New()
	server.Use(echoMiddleware.Recover())
	server.Use(middleware.TraceMiddleware)

	// init health
	health.HealthController(server.Group("/api/health"))

	// init swagger
	server.GET("/swagger/*", echoSwagger.WrapHandler)

	// init v1api
	v1api.New(server.Group("/api/v1"), v1api.Depends{
		Cfg:          cfg,
		OrderService: service.OrderService,
	})

	server.HideBanner = true
	server.HidePort = true

	return &Server{
		server: server,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	logger.Info("http server started", zap.String("port", s.cfg.Port))
	return s.server.Start(":" + s.cfg.Port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
