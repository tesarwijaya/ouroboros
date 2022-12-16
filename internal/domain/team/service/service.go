package service

import (
	"context"

	player_repository "github.com/tesarwijaya/ouroboros/internal/domain/player/repository"
	"github.com/tesarwijaya/ouroboros/internal/domain/team/model"
	"github.com/tesarwijaya/ouroboros/internal/domain/team/repository"
	"go.uber.org/dig"
)

type TeamService interface {
	FindAll(ctx context.Context) ([]model.TeamModel, error)
	FindByID(ctx context.Context, id int64) (model.TeamModel, error)
	FindTeamPlayer(ctx context.Context, id int64) (model.TeamPlayerRespModel, error)
	Insert(ctx context.Context, payload model.TeamModel) (model.TeamModel, error)
}

type TeamServiceImpl struct {
	dig.In
	Repo       repository.TeamRepository
	PlayerRepo player_repository.PlayerRepository
}

func NewTeamService(svc TeamServiceImpl) TeamService {
	return &svc
}

func (s *TeamServiceImpl) FindAll(ctx context.Context) ([]model.TeamModel, error) {
	return s.Repo.FindAll(ctx)
}

func (s *TeamServiceImpl) FindByID(ctx context.Context, id int64) (model.TeamModel, error) {
	return s.Repo.FindByID(ctx, id)
}

func (s *TeamServiceImpl) Insert(ctx context.Context, payload model.TeamModel) (model.TeamModel, error) {
	if err := s.Repo.Insert(ctx, payload); err != nil {
		return model.TeamModel{}, err
	}

	return payload, nil
}

func (s *TeamServiceImpl) FindTeamPlayer(ctx context.Context, id int64) (model.TeamPlayerRespModel, error) {
	team, err := s.FindByID(ctx, id)
	if err != nil {
		return model.TeamPlayerRespModel{}, err
	}

	players, _ := s.PlayerRepo.FindByTeamID(ctx, id)
	return model.TeamPlayerRespModel{
		TeamModel: team,
		Players:   players,
	}, nil
}
