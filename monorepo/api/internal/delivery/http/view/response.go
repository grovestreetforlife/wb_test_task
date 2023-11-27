package view

import (
	"errors"
	"github.com/dany-ykl/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"wb_test_task/api/internal/common"
	"wb_test_task/api/internal/domain"
)

type Response struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Body   any    `json:"body"`
	Error  string `json:"error"`
}

func SuccessResponse(c echo.Context, code int, body any) error {
	return c.JSON(code, Response{
		Code:   http.StatusText(code),
		Status: "ok",
		Body:   body,
	})
}

func ErrorResponse(c echo.Context, err error) error {
	var wrapError common.WrapError
	if !errors.As(err, &wrapError) {
		logger.Warn("error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Response{
			Code:   http.StatusText(http.StatusInternalServerError),
			Status: "fail",
			Body:   nil,
			Error:  err.Error(),
		})
	} else {
		logger.Warn("error", zap.String("msg", wrapError.Msg), zap.Error(wrapError.Err))
		return c.JSON(wrapError.Code, Response{
			Code:   http.StatusText(wrapError.Code),
			Status: "fail",
			Body:   wrapError.Body,
			Error:  wrapError.Msg,
		})
	}
}

func ErrorResponseSwitch(c echo.Context, err error) error {
	var httpErr common.WrapError
	if !errors.As(err, &httpErr) {
		return ErrorResponse(c, common.WrapError{Code: http.StatusInternalServerError, Err: err, Msg: err.Error()})
	}

	if httpErr.Code != 0 {
		return ErrorResponse(c, err)
	}

	switch httpErr.Err {
	case domain.ErrOrderNotExists, domain.ErrItemsNotExists:
		return ErrorResponse(c, common.WrapError{Code: http.StatusNotFound, Err: httpErr.Err, Msg: httpErr.Msg})
	case ErrInvalidOrderID, domain.ErrInvalidSyntax:
		return ErrorResponse(c, common.WrapError{Code: http.StatusBadRequest, Err: httpErr.Err, Msg: httpErr.Msg})
	default:
		return ErrorResponse(c, common.WrapError{Code: http.StatusInternalServerError, Err: httpErr.Err, Msg: httpErr.Msg})
	}
}
