package middleware

import (
	"github.com/dany-ykl/tracer"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TraceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ctx, span := tracer.StartTrace(ctx, "http-request")
		defer span.End()

		span.SetStatus(codes.Ok, "http request")
		span.SetAttributes(attribute.String("Method", c.Request().Method))
		span.SetAttributes(attribute.String("URL", c.Request().URL.String()))

		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
