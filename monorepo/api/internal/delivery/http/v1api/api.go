package v1api

import (
	"github.com/labstack/echo/v4"
	"wb_test_task/api/internal/config"
)

type API struct {
	cfg          config.HttpServer
	orderService orderService
}

type Depends struct {
	Cfg          config.HttpServer
	OrderService orderService
}

func New(group *echo.Group, depends Depends) *API {
	api := &API{
		cfg:          depends.Cfg,
		orderService: depends.OrderService,
	}
	api.initControllers(group)
	return api
}

// initControllers инициализация контроллеров
func (a *API) initControllers(group *echo.Group) {
	a.orderController(group.Group("/orders"))
}
