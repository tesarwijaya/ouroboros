package repository

import (
	"context"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/tesarwijaya/ouroboros/internal/domain/event/model"
	"go.uber.org/dig"
)

type EventRepository interface {
	Insert(ctx context.Context, payload model.Event) error
}

type EventRepositoryImpl struct {
	dig.In
	Db *esdb.Client
}

func NewTeamReposity(repo EventRepositoryImpl) EventRepository {
	return &repo
}

func (r *EventRepositoryImpl) Insert(ctx context.Context, payload model.Event) error {
	eventData := esdb.EventData{
		EventID:     payload.ID,
		EventType:   payload.Type,
		ContentType: payload.ContentType,
		Data:        payload.Data,
		Metadata:    payload.Metadata,
	}

	_, err := r.Db.AppendToStream(ctx, payload.ID.String(), esdb.AppendToStreamOptions{}, eventData)
	if err != nil {
		return err
	}

	return nil
}
