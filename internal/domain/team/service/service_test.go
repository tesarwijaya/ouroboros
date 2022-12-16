package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	player_model "github.com/tesarwijaya/ouroboros/internal/domain/player/model"
	player_repository "github.com/tesarwijaya/ouroboros/internal/domain/player/repository"
	"github.com/tesarwijaya/ouroboros/internal/domain/team/model"
	"github.com/tesarwijaya/ouroboros/internal/domain/team/repository"
	"github.com/tesarwijaya/ouroboros/internal/domain/team/service"
)

type resolverFn func(repo *repository.MockTeamRepository, playerRepo *player_repository.MockPlayerRepository)

func createService(t *testing.T, resolver resolverFn) (*service.TeamServiceImpl, gomock.Controller) {
	ctrl := gomock.NewController(t)

	repo := repository.NewMockTeamRepository(ctrl)
	playerRepo := player_repository.NewMockPlayerRepository(ctrl)
	resolver(repo, playerRepo)

	return &service.TeamServiceImpl{
		Repo:       repo,
		PlayerRepo: playerRepo,
	}, *ctrl
}

func Test_NewTeamService(t *testing.T) {
	svc := service.NewTeamService(service.TeamServiceImpl{})

	assert.Implements(t, (*service.TeamService)(nil), svc)
}

func Test_FindAll(t *testing.T) {
	testCases := []struct {
		Name      string
		Resolver  resolverFn
		Expect    []model.TeamModel
		ExpectErr error
	}{
		{
			Name: "when_success",
			Resolver: func(repo *repository.MockTeamRepository, playerRepo *player_repository.MockPlayerRepository) {
				repo.EXPECT().FindAll(gomock.Any()).
					Return([]model.TeamModel{{Name: "some-team-name"}}, nil)
			},
			Expect: []model.TeamModel{{Name: "some-team-name"}},
		},
		{
			Name: "when_not_success",
			Resolver: func(repo *repository.MockTeamRepository, playerRepo *player_repository.MockPlayerRepository) {
				repo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("some-error"))
			},
			ExpectErr: errors.New("some-error"),
		},
	}

	for _, test := range testCases {
		svc, mock := createService(t, test.Resolver)
		defer mock.Finish()

		actual, err := svc.FindAll(context.Background())

		if test.ExpectErr == nil {
			assert.Equal(t, test.Expect, actual)
			assert.Nil(t, err)
		}

		assert.Equal(t, test.ExpectErr, err)
	}
}

func Test_FindByID(t *testing.T) {
	testCases := []struct {
		Name      string
		Param     int64
		Resolver  resolverFn
		Expect    model.TeamModel
		ExpectErr error
	}{
		{
			Name:  "when_success",
			Param: 1,
			Resolver: func(repo *repository.MockTeamRepository, playerRepo *player_repository.MockPlayerRepository) {
				repo.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(model.TeamModel{Name: "some-team-name"}, nil)
			},
			Expect: model.TeamModel{Name: "some-team-name"},
		},
		{
			Name:  "when_not_success",
			Param: 1,
			Resolver: func(repo *repository.MockTeamRepository, playerRepo *player_repository.MockPlayerRepository) {
				repo.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(model.TeamModel{}, errors.New("some-error"))
			},
			ExpectErr: errors.New("some-error"),
		},
	}

	for _, test := range testCases {
		svc, mock := createService(t, test.Resolver)
		defer mock.Finish()

		actual, err := svc.FindByID(context.Background(), test.Param)

		if test.ExpectErr == nil {
			assert.Equal(t, test.Expect, actual)
			assert.Nil(t, err)
		}

		assert.Equal(t, test.ExpectErr, err)
	}
}

func Test_Insert(t *testing.T) {
	testCases := []struct {
		Name      string
		Param     model.TeamModel
		Resolver  resolverFn
		Expect    model.TeamModel
		ExpectErr error
	}{
		{
			Name:  "when_success",
			Param: model.TeamModel{Name: "some-team-name"},
			Resolver: func(repo *repository.MockTeamRepository, playerRepo *player_repository.MockPlayerRepository) {
				repo.EXPECT().Insert(gomock.Any(), model.TeamModel{Name: "some-team-name"}).
					Return(nil)
			},
			Expect: model.TeamModel{Name: "some-team-name"},
		},
		{
			Name:  "when_not_success",
			Param: model.TeamModel{Name: "some-team-name"},
			Resolver: func(repo *repository.MockTeamRepository, playerRepo *player_repository.MockPlayerRepository) {
				repo.EXPECT().Insert(gomock.Any(), model.TeamModel{Name: "some-team-name"}).
					Return(errors.New("some-error"))
			},
			ExpectErr: errors.New("some-error"),
		},
	}

	for _, test := range testCases {
		svc, mock := createService(t, test.Resolver)
		defer mock.Finish()

		actual, err := svc.Insert(context.Background(), test.Param)

		if test.ExpectErr == nil {
			assert.Equal(t, test.Expect, actual)
			assert.Nil(t, err)
		}

		assert.Equal(t, test.ExpectErr, err)
	}
}

func Test_FindTeamPlayer(t *testing.T) {
	testCases := []struct {
		Name      string
		Param     int64
		Resolver  resolverFn
		Expect    model.TeamPlayerRespModel
		ExpectErr error
	}{
		{
			Name:  "when_success",
			Param: 1,
			Resolver: func(repo *repository.MockTeamRepository, playerRepo *player_repository.MockPlayerRepository) {
				repo.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(model.TeamModel{Name: "some-team-name"}, nil)
				playerRepo.EXPECT().FindByTeamID(gomock.Any(), int64(1)).
					Return([]player_model.PlayerModel{{Name: "some-player", TeamID: 1}}, nil)

			},
			Expect: model.TeamPlayerRespModel{
				TeamModel: model.TeamModel{Name: "some-team-name"},
				Players:   []player_model.PlayerModel{{Name: "some-player", TeamID: 1}},
			},
		},
		{
			Name:  "when_not_success",
			Param: 1,
			Resolver: func(repo *repository.MockTeamRepository, playerRepo *player_repository.MockPlayerRepository) {
				repo.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(model.TeamModel{}, errors.New("some-error"))
			},
			ExpectErr: errors.New("some-error"),
		},
	}

	for _, test := range testCases {
		svc, mock := createService(t, test.Resolver)
		defer mock.Finish()

		actual, err := svc.FindTeamPlayer(context.Background(), test.Param)

		if test.ExpectErr == nil {
			assert.Equal(t, test.Expect, actual)
			assert.Nil(t, err)
		}

		assert.Equal(t, test.ExpectErr, err)
	}
}
