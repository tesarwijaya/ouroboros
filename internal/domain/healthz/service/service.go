package service

import (
	"context"
	"database/sql"

	"go.uber.org/dig"
)

type HealthzService interface {
	Healthz(ctx context.Context) (map[string]interface{}, error)
}

type HealthzServiceImpl struct {
	dig.In
	Sql *sql.DB
}

func NewHealthzService(svc HealthzServiceImpl) HealthzService {
	return &svc
}

func (s *HealthzServiceImpl) Healthz(ctx context.Context) (map[string]interface{}, error) {
	DBStatus := "UP!"
	err := s.Sql.Ping()
	if err != nil {
		DBStatus = err.Error()
	}

	return map[string]interface{}{
		"DBStatus": DBStatus,
	}, nil
}
