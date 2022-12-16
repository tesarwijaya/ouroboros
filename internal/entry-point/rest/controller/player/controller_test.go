package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/model"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/service"
	controller "github.com/tesarwijaya/ouroboros/internal/entry-point/rest/controller/player"
)

type ResolverFn func(svc *service.MockPlayerService)

func createController(t *testing.T, resolver ResolverFn) (controller.PlayerController, gomock.Controller) {
	ctrl := gomock.NewController(t)

	svc := service.NewMockPlayerService(ctrl)
	resolver(svc)

	return controller.PlayerController{
		Service: svc,
	}, *ctrl
}

func Test_FindAll(t *testing.T) {
	testCases := []struct {
		Name             string
		Resolver         ResolverFn
		ExpectBody       string
		ExpectStatusCode int64
		ExpectErr        error
	}{
		{
			Name: "when_success",
			Resolver: func(svc *service.MockPlayerService) {
				svc.EXPECT().FindAll(gomock.Any()).
					Return([]model.PlayerModel{{Name: "some-player-name"}}, nil)
			},
			ExpectStatusCode: 200,
			ExpectBody:       "[{\"id\":0,\"name\":\"some-player-name\"}]\n",
		},
		{
			Name: "when_not_success",
			Resolver: func(svc *service.MockPlayerService) {
				svc.EXPECT().FindAll(gomock.Any()).
					Return([]model.PlayerModel{}, errors.New("some-error"))
			},
			ExpectStatusCode: 500,
			ExpectErr:        echo.NewHTTPError(http.StatusInternalServerError, "some-error"),
		},
	}

	for _, test := range testCases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/player", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		controller, mock := createController(t, test.Resolver)
		defer mock.Finish()

		err := controller.FindAll(c)
		if test.ExpectErr == nil {
			assert.Equal(t, test.ExpectBody, rec.Body.String())
			assert.Equal(t, http.StatusOK, rec.Code)
		} else {
			assert.Equal(t, test.ExpectErr, err)
		}
	}
}

func Test_FindByID(t *testing.T) {
	testCases := []struct {
		Name             string
		QueryString      string
		Resolver         ResolverFn
		ExpectBody       string
		ExpectStatusCode int64
		ExpectErr        error
	}{
		{
			Name:        "when_success",
			QueryString: "1",
			Resolver: func(svc *service.MockPlayerService) {
				svc.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(model.PlayerModel{Name: "some-player-name"}, nil)
			},
			ExpectStatusCode: 200,
			ExpectBody:       "{\"id\":0,\"name\":\"some-player-name\"}\n",
		},
	}

	for _, test := range testCases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/player/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(test.QueryString)

		controller, mock := createController(t, test.Resolver)
		defer mock.Finish()

		err := controller.FindByID(c)
		if test.ExpectErr == nil {
			assert.Equal(t, test.ExpectBody, rec.Body.String())
			assert.Equal(t, http.StatusOK, rec.Code)
		} else {
			assert.Equal(t, test.ExpectErr, err)
		}
	}
}

func Test_Insert(t *testing.T) {
	testCases := []struct {
		Name             string
		Body             model.PlayerModel
		Resolver         ResolverFn
		ExpectBody       string
		ExpectStatusCode int64
		ExpectErr        error
	}{
		{
			Name: "when_success",
			Body: model.PlayerModel{Name: "some-player-name"},
			Resolver: func(svc *service.MockPlayerService) {
				svc.EXPECT().Insert(gomock.Any(), model.PlayerModel{Name: "some-player-name"}).
					Return(model.PlayerModel{Name: "some-player-name"}, nil)
			},
			ExpectStatusCode: 200,
			ExpectBody:       "{\"id\":0,\"name\":\"some-player-name\"}\n",
		},
	}

	for _, test := range testCases {
		e := echo.New()
		var body bytes.Buffer
		_ = json.NewEncoder(&body).Encode(test.Body)

		req := httptest.NewRequest(http.MethodPost, "/player", &body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		controller, mock := createController(t, test.Resolver)
		defer mock.Finish()

		err := controller.Insert(c)
		if test.ExpectErr == nil {
			assert.Equal(t, test.ExpectBody, rec.Body.String())
			assert.Equal(t, http.StatusOK, rec.Code)
		} else {
			assert.Equal(t, test.ExpectErr, err)
		}
	}
}
