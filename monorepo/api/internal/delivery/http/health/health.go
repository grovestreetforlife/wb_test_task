package health

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HealthController(g *echo.Group) {
	g.GET("/live", func(c echo.Context) error {
		return c.String(http.StatusOK, "Health")
	})

	g.GET("/readiness", func(c echo.Context) error {
		return c.String(http.StatusOK, "Readiness")
	})
}
