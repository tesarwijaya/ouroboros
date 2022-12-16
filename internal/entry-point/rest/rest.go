package rest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/tesarwijaya/ouroboros/docs"
	"github.com/tesarwijaya/ouroboros/internal/config"
	healthz_controller "github.com/tesarwijaya/ouroboros/internal/entry-point/rest/controller/healthz"
	player_controller "github.com/tesarwijaya/ouroboros/internal/entry-point/rest/controller/player"
	team_controller "github.com/tesarwijaya/ouroboros/internal/entry-point/rest/controller/team"
	"go.uber.org/dig"
)

type RestController struct {
	dig.In
	HealthzController healthz_controller.HealthzController
	PlayerController  player_controller.PlayerController
	TeamController    team_controller.TeamController
}

type RestServer struct {
	Server *echo.Echo
	Config *config.Config
}

// @title          Night owl API
// @version        1.0
// @description    Night owl API documentation.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8000
// @BasePath /
func NewRestServer(c *config.Config, controllers RestController) RestServer {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	controllers.HealthzController.SetRouter(e)
	controllers.PlayerController.SetRouter(e)
	controllers.TeamController.SetRouter(e)

	return RestServer{
		Server: e,
		Config: c,
	}
}

func (s *RestServer) Start() error {
	return s.Server.Start(fmt.Sprintf(":%s", s.Config.Port))
}
