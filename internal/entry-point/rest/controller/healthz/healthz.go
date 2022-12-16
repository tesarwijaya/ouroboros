package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tesarwijaya/ouroboros/internal/domain/healthz/service"
)

type HealthzController struct {
	Service service.HealthzService
}

func NewHealthzController(service service.HealthzService) HealthzController {
	return HealthzController{
		Service: service,
	}
}

func (c *HealthzController) SetRouter(ec *echo.Echo) {
	ec.GET("/healthz", c.Healthz)
}

func (c *HealthzController) Healthz(ec echo.Context) error {
	res, err := c.Service.Healthz(ec.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusOK, res)
}
