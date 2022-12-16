package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/model"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/service"
)

type PlayerController struct {
	Service service.PlayerService
}

func NewPlayerController(service service.PlayerService) PlayerController {
	return PlayerController{
		Service: service,
	}
}

func (c *PlayerController) SetRouter(ec *echo.Echo) {
	ec.GET("/player", c.FindAll)
	ec.GET("/player/:id", c.FindByID)
	ec.POST("/player", c.Insert)
	ec.PATCH("/player/transfer", c.Transfer)
}

// FindAll godoc
// @Summary      Show all player
// @Description  get all player
// @Tags         Player
// @Accept       json
// @Produce      json
// @Success      200  {object}  []model.PlayerModel
// @Failure      400  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /player [get]
func (c *PlayerController) FindAll(ec echo.Context) error {
	res, err := c.Service.FindAll(ec.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusOK, res)
}

// FindByID godoc
// @Summary      Get player by id
// @Description  get player by id
// @Tags         Player
// @Accept       json
// @Produce      json
// @param        id path int true "player id"
// @Success      200  {object}  model.PlayerModel
// @Failure      400  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /player/{id} [get]
func (c *PlayerController) FindByID(ec echo.Context) error {
	idParam := ec.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := c.Service.FindByID(ec.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusOK, res)
}

// Insert godoc
// @Summary      Insert player
// @Description  insert a player with team id
// @Tags         Player
// @Accept       json
// @Produce      json
// @param        id body model.PlayerModel true "body"
// @Success      200  {object}  model.PlayerModel
// @Failure      400  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /player [post]
func (c *PlayerController) Insert(ec echo.Context) error {
	var payload model.PlayerModel

	if err := ec.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := c.Service.Insert(ec.Request().Context(), payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusOK, res)
}

// Transfer godoc
// @Summary      Transfer player
// @Description  Transfer a player to team id
// @Tags         Player
// @Accept       json
// @Produce      json
// @param        id body service.TransferPayload true "body"
// @Failure      400  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /player/transfer [post]
func (c *PlayerController) Transfer(ec echo.Context) error {
	var payload service.TransferPayload

	if err := ec.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := c.Service.Transfer(ec.Request().Context(), payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusNoContent, nil)
}
