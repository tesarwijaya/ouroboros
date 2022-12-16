package repository

import (
	"context"
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
	"github.com/tesarwijaya/ouroboros/internal/domain/team/model"
	"go.uber.org/dig"
)

const (
	TEAM_TABLE_NAME = "team"
)

type TeamRepository interface {
	FindAll(ctx context.Context) ([]model.TeamModel, error)
	FindByID(ctx context.Context, id int64) (model.TeamModel, error)
	Insert(ctx context.Context, payload model.TeamModel) error
}

type TeamRepositoryImpl struct {
	dig.In
	Db *sql.DB
}

func NewTeamReposity(repo TeamRepositoryImpl) TeamRepository {
	return &repo
}

func (r *TeamRepositoryImpl) FindAll(ctx context.Context) ([]model.TeamModel, error) {
	var res []model.TeamModel
	q := sqlbuilder.NewSelectBuilder()
	query, _ := q.Select("*").From(TEAM_TABLE_NAME).BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.Db.Query(query)
	if err != nil {
		return []model.TeamModel{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.TeamModel

		if err = rows.Scan(
			&item.ID,
			&item.Name,
		); err != nil {
			return []model.TeamModel{}, err
		}

		res = append(res, item)
	}

	if err = rows.Err(); err != nil {
		return []model.TeamModel{}, err
	}

	return res, nil
}

func (r *TeamRepositoryImpl) FindByID(ctx context.Context, id int64) (model.TeamModel, error) {
	var res model.TeamModel
	q := sqlbuilder.NewSelectBuilder()
	query, args := q.Select("*").From(TEAM_TABLE_NAME).Where(q.Equal("id", id)).BuildWithFlavor(sqlbuilder.PostgreSQL)

	row := r.Db.QueryRow(query, args...)
	if err := row.Err(); err != nil {
		return model.TeamModel{}, err
	}

	if err := row.Scan(
		&res.ID,
		&res.Name,
	); err != nil {
		return model.TeamModel{}, err
	}

	return res, nil
}

func (r *TeamRepositoryImpl) Insert(ctx context.Context, payload model.TeamModel) error {
	q := sqlbuilder.NewInsertBuilder()

	query, args := q.InsertInto(TEAM_TABLE_NAME).Cols("name").Values(payload.Name).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := r.Db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
