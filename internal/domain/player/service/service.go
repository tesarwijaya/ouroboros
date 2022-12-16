package service

import (
	"context"
	"encoding/json"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"
	event_model "github.com/tesarwijaya/ouroboros/internal/domain/event/model"
	event_repository "github.com/tesarwijaya/ouroboros/internal/domain/event/repository"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/model"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/repository"
	team_repository "github.com/tesarwijaya/ouroboros/internal/domain/team/repository"
	"go.uber.org/dig"
)

type (
	TransferPayload struct {
		PlayerID int64
		TeamID   int64
	}
)

type PlayerService interface {
	FindAll(ctx context.Context) ([]model.PlayerModel, error)
	FindByID(ctx context.Context, id int64) (model.PlayerModel, error)
	Insert(ctx context.Context, payload model.PlayerModel) (model.PlayerModel, error)
	Transfer(ctx context.Context, payload TransferPayload) error
}

type PlayerServiceImpl struct {
	dig.In
	Repo      repository.PlayerRepository
	TeamRepo  team_repository.TeamRepository
	EventRepo event_repository.EventRepository
}

func NewPlayerService(svc PlayerServiceImpl) PlayerService {
	return &svc
}

func (s *PlayerServiceImpl) FindAll(ctx context.Context) ([]model.PlayerModel, error) {
	return s.Repo.FindAll(ctx)
}

func (s *PlayerServiceImpl) FindByID(ctx context.Context, id int64) (model.PlayerModel, error) {
	return s.Repo.FindByID(ctx, id)
}

func (s *PlayerServiceImpl) Insert(ctx context.Context, payload model.PlayerModel) (model.PlayerModel, error) {
	_, err := s.TeamRepo.FindByID(ctx, payload.TeamID)
	if err != nil {
		return model.PlayerModel{}, err
	}

	if err := s.Repo.Insert(ctx, payload); err != nil {
		return model.PlayerModel{}, err
	}

	return payload, nil
}

func (s *PlayerServiceImpl) Transfer(ctx context.Context, payload TransferPayload) error {
	gen := uuid.NewGen()

	payloadByte, _ := json.Marshal(payload)

	currPlayer, err := s.Repo.FindByID(ctx, payload.PlayerID)
	if err != nil {
		return err
	}

	playerOutByte, _ := json.Marshal(TransferPayload{
		PlayerID: currPlayer.ID,
		TeamID:   currPlayer.TeamID,
	})

	outId, _ := gen.NewV4()
	err = s.EventRepo.Insert(ctx, event_model.Event{
		ID:          outId,
		Type:        "player_transfer_out",
		ContentType: esdb.JsonContentType,
		Data:        playerOutByte,
	})
	if err != nil {
		return err
	}

	inId, _ := gen.NewV4()
	err = s.EventRepo.Insert(ctx, event_model.Event{
		ID:          inId,
		Type:        "player_transfer_in",
		ContentType: esdb.JsonContentType,
		Data:        payloadByte,
	})
	if err != nil {
		return err
	}

	return nil
}
