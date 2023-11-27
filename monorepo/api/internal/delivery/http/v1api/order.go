package v1api

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"wb_test_task/api/internal/common"
	"wb_test_task/api/internal/delivery/http/view"
	"wb_test_task/libs/model"
)

//go:generate mockgen -source=order.go -destination=mocks/mock.go
type orderService interface {
	GetByID(ctx context.Context, id string) (*model.Order, error)
}

func (a *API) orderController(g *echo.Group) {
	g.GET("/:id", func(c echo.Context) error {
		return a.getOrder(c)
	})
}

func (a *API) getOrder(c echo.Context) error {
	id := c.Param("id")
	if err := validateOrderID(id); err != nil {
		return view.ErrorResponse(c, err)
	}

	order, err := a.orderService.GetByID(c.Request().Context(), id)
	if err != nil {
		return view.ErrorResponseSwitch(c, err)
	}

	return view.SuccessResponse(c, http.StatusOK, order)
}

func validateOrderID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return common.WrapError{Code: http.StatusBadRequest, Err: view.ErrInvalidOrderID, Msg: "invalid order uuid id"}
	}

	return nil
}
