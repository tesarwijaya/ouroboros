package repository

import (
	"context"
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/model"
	"go.uber.org/dig"
)

const (
	PLAYER_TABLE_NAME = "player"
)

type PlayerRepository interface {
	FindAll(ctx context.Context) ([]model.PlayerModel, error)
	FindByID(ctx context.Context, id int64) (model.PlayerModel, error)
	FindByTeamID(ctx context.Context, teamID int64) ([]model.PlayerModel, error)
	Insert(ctx context.Context, payload model.PlayerModel) error
}

type PlayerRepositoryImpl struct {
	dig.In
	Db *sql.DB
}

func NewPlayerReposity(repo PlayerRepositoryImpl) PlayerRepository {
	return &repo
}

func (r *PlayerRepositoryImpl) FindAll(ctx context.Context) ([]model.PlayerModel, error) {
	var res []model.PlayerModel
	q := sqlbuilder.NewSelectBuilder()
	query, _ := q.Select("*").From(PLAYER_TABLE_NAME).Build()

	rows, err := r.Db.Query(query)
	if err != nil {
		return []model.PlayerModel{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.PlayerModel

		if err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.TeamID,
		); err != nil {
			return []model.PlayerModel{}, err
		}

		res = append(res, item)
	}

	err = rows.Err()
	if err != nil {
		return []model.PlayerModel{}, err
	}

	return res, nil
}

func (r *PlayerRepositoryImpl) FindByID(ctx context.Context, id int64) (model.PlayerModel, error) {
	var res model.PlayerModel
	q := sqlbuilder.NewSelectBuilder()
	query, args := q.Select("*").From(PLAYER_TABLE_NAME).Where(q.Equal("id", id)).BuildWithFlavor(sqlbuilder.PostgreSQL)

	row := r.Db.QueryRow(query, args...)
	if err := row.Err(); err != nil {
		return model.PlayerModel{}, err
	}

	if err := row.Scan(
		&res.ID,
		&res.Name,
		&res.TeamID,
	); err != nil {
		return model.PlayerModel{}, err
	}

	return res, nil
}

func (r *PlayerRepositoryImpl) FindByTeamID(ctx context.Context, teamID int64) ([]model.PlayerModel, error) {
	var res []model.PlayerModel
	q := sqlbuilder.NewSelectBuilder()
	query, args := q.Select("*").From(PLAYER_TABLE_NAME).Where(q.Equal("team_id", teamID)).BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.Db.Query(query, args...)
	if err != nil {
		return []model.PlayerModel{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.PlayerModel
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.TeamID,
		); err != nil {
			return []model.PlayerModel{}, err
		}

		res = append(res, item)
	}

	if err = rows.Err(); err != nil {
		return []model.PlayerModel{}, err
	}

	return res, nil
}

func (r *PlayerRepositoryImpl) Insert(ctx context.Context, payload model.PlayerModel) error {
	q := sqlbuilder.NewInsertBuilder()
	query, args := q.InsertInto(PLAYER_TABLE_NAME).
		Cols("name", "team_id").
		Values(payload.Name, payload.TeamID).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := r.Db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
