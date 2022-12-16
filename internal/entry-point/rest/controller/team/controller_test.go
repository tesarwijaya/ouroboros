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
	"github.com/tesarwijaya/ouroboros/internal/domain/team/model"
	"github.com/tesarwijaya/ouroboros/internal/domain/team/service"
	controller "github.com/tesarwijaya/ouroboros/internal/entry-point/rest/controller/team"
)

type ResolverFn func(svc *service.MockTeamService)

func createController(t *testing.T, resolver ResolverFn) (controller.TeamController, gomock.Controller) {
	ctrl := gomock.NewController(t)

	svc := service.NewMockTeamService(ctrl)
	resolver(svc)

	return controller.TeamController{
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
			Resolver: func(svc *service.MockTeamService) {
				svc.EXPECT().FindAll(gomock.Any()).
					Return([]model.TeamModel{{Name: "some-team-name"}}, nil)
			},
			ExpectStatusCode: 200,
			ExpectBody:       "[{\"name\":\"some-team-name\"}]\n",
		},
		{
			Name: "when_not_success",
			Resolver: func(svc *service.MockTeamService) {
				svc.EXPECT().FindAll(gomock.Any()).
					Return([]model.TeamModel{}, errors.New("some-error"))
			},
			ExpectStatusCode: 500,
			ExpectErr:        echo.NewHTTPError(http.StatusInternalServerError, "some-error"),
		},
	}

	for _, test := range testCases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/team", nil)
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
			Resolver: func(svc *service.MockTeamService) {
				svc.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(model.TeamModel{Name: "some-team-name"}, nil)
			},
			ExpectStatusCode: 200,
			ExpectBody:       "{\"name\":\"some-team-name\"}\n",
		},
	}

	for _, test := range testCases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/team/:id", nil)
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
		Body             model.TeamModel
		Resolver         ResolverFn
		ExpectBody       string
		ExpectStatusCode int64
		ExpectErr        error
	}{
		{
			Name: "when_success",
			Body: model.TeamModel{Name: "some-team-name"},
			Resolver: func(svc *service.MockTeamService) {
				svc.EXPECT().Insert(gomock.Any(), model.TeamModel{Name: "some-team-name"}).
					Return(model.TeamModel{Name: "some-team-name"}, nil)
			},
			ExpectStatusCode: 200,
			ExpectBody:       "{\"name\":\"some-team-name\"}\n",
		},
	}

	for _, test := range testCases {
		e := echo.New()
		var body bytes.Buffer
		_ = json.NewEncoder(&body).Encode(test.Body)

		req := httptest.NewRequest(http.MethodPost, "/team", &body)
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
