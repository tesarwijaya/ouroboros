package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/model"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/repository"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/service"
	team_model "github.com/tesarwijaya/ouroboros/internal/domain/team/model"
	team_repository "github.com/tesarwijaya/ouroboros/internal/domain/team/repository"
)

type resolverFn func(repo *repository.MockPlayerRepository, teamRepo *team_repository.MockTeamRepository)

func createService(t *testing.T, resolver resolverFn) (*service.PlayerServiceImpl, gomock.Controller) {
	ctrl := gomock.NewController(t)

	repo := repository.NewMockPlayerRepository(ctrl)
	teamRepo := team_repository.NewMockTeamRepository(ctrl)
	resolver(repo, teamRepo)

	return &service.PlayerServiceImpl{
		Repo:     repo,
		TeamRepo: teamRepo,
	}, *ctrl
}

func Test_NewPlayerService(t *testing.T) {
	svc := service.NewPlayerService(service.PlayerServiceImpl{})

	assert.Implements(t, (*service.PlayerService)(nil), svc)
}

func Test_FindAll(t *testing.T) {
	testCases := []struct {
		Name      string
		Resolver  resolverFn
		Expect    []model.PlayerModel
		ExpectErr error
	}{
		{
			Name: "when_success",
			Resolver: func(repo *repository.MockPlayerRepository, teamRepo *team_repository.MockTeamRepository) {
				repo.EXPECT().FindAll(gomock.Any()).
					Return([]model.PlayerModel{{Name: "some-player-name"}}, nil)
			},
			Expect: []model.PlayerModel{{Name: "some-player-name"}},
		},
		{
			Name: "when not success",
			Resolver: func(repo *repository.MockPlayerRepository, teamRepo *team_repository.MockTeamRepository) {
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
		Expect    model.PlayerModel
		ExpectErr error
	}{
		{
			Name:  "when_success",
			Param: 1,
			Resolver: func(repo *repository.MockPlayerRepository, teamRepo *team_repository.MockTeamRepository) {
				repo.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(model.PlayerModel{Name: "some-player-name"}, nil)
			},
			Expect: model.PlayerModel{Name: "some-player-name"},
		},
		{
			Name:  "when_not_success",
			Param: 1,
			Resolver: func(repo *repository.MockPlayerRepository, teamRepo *team_repository.MockTeamRepository) {
				repo.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(model.PlayerModel{}, errors.New("some-error"))
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
		Param     model.PlayerModel
		Resolver  resolverFn
		Expect    model.PlayerModel
		ExpectErr error
	}{
		{
			Name:  "when_success",
			Param: model.PlayerModel{Name: "some-player-name", TeamID: 1},
			Resolver: func(repo *repository.MockPlayerRepository, teamRepo *team_repository.MockTeamRepository) {
				teamRepo.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(team_model.TeamModel{}, nil)
				repo.EXPECT().Insert(gomock.Any(), model.PlayerModel{Name: "some-player-name", TeamID: 1}).
					Return(nil)
			},
			Expect: model.PlayerModel{Name: "some-player-name", TeamID: 1},
		},
		{
			Name:  "when_not_success",
			Param: model.PlayerModel{Name: "some-player-name", TeamID: 1},
			Resolver: func(repo *repository.MockPlayerRepository, teamRepo *team_repository.MockTeamRepository) {
				teamRepo.EXPECT().FindByID(gomock.Any(), int64(1)).
					Return(team_model.TeamModel{}, nil)
				repo.EXPECT().Insert(gomock.Any(), model.PlayerModel{Name: "some-player-name", TeamID: 1}).
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
