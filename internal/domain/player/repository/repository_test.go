package repository_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/model"
	"github.com/tesarwijaya/ouroboros/internal/domain/player/repository"
)

type mockFn func(db sqlmock.Sqlmock)

func createRepo(mockFn mockFn) repository.PlayerRepository {
	db, mock, _ := sqlmock.New()

	mockFn(mock)
	repo := repository.NewPlayerReposity(repository.PlayerRepositoryImpl{
		Db: db,
	})

	return repo
}

func Test_FindAll(t *testing.T) {
	testCases := []struct {
		Name        string
		MockFn      mockFn
		Expected    []model.PlayerModel
		ExpectedErr string
	}{
		{
			Name: "when success",
			MockFn: func(db sqlmock.Sqlmock) {
				db.ExpectQuery(regexp.QuoteMeta("SELECT * FROM player")).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "team_id"}).
							AddRow(int64(1), "some-player-name", int64(1)),
					)
			},
			Expected: []model.PlayerModel{{
				ID:     1,
				Name:   "some-player-name",
				TeamID: 1,
			}},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			repo := createRepo(test.MockFn)

			actual, err := repo.FindAll(context.Background())
			if test.ExpectedErr != "" {
				assert.EqualError(t, err, test.ExpectedErr)
			} else {
				assert.Equal(t, test.Expected, actual)
				assert.Nil(t, err)
			}
		})
	}
}

func Test_FindByID(t *testing.T) {
	testCases := []struct {
		Name      string
		Param     int64
		mockFn    mockFn
		Expect    model.PlayerModel
		ExpectErr error
	}{
		{
			Name:  "when success",
			Param: 1,
			mockFn: func(db sqlmock.Sqlmock) {
				db.ExpectQuery(regexp.QuoteMeta("SELECT * FROM player WHERE id = $1")).
					WithArgs(int64(1)).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "team_id"}).
							AddRow(int64(1), "some-player-name", int64(1)),
					)
			},
			Expect: model.PlayerModel{
				ID:     1,
				Name:   "some-player-name",
				TeamID: 1,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			repo := createRepo(test.mockFn)

			actual, err := repo.FindByID(context.Background(), test.Param)
			if test.ExpectErr == nil {
				assert.Equal(t, test.Expect, actual)
				assert.Nil(t, err)
			}

			assert.Equal(t, test.ExpectErr, err)
		})
	}
}

func Test_FindByTeamID(t *testing.T) {
	testCases := []struct {
		Name      string
		Param     int64
		mockFn    mockFn
		Expect    []model.PlayerModel
		ExpectErr error
	}{
		{
			Name:  "when success",
			Param: 1,
			mockFn: func(db sqlmock.Sqlmock) {
				db.ExpectQuery(regexp.QuoteMeta("SELECT * FROM player WHERE team_id = $1")).
					WithArgs(int64(1)).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "team_id"}).
							AddRow(1, "some-player-name", 1),
					)
			},
			Expect: []model.PlayerModel{{
				ID:     1,
				Name:   "some-player-name",
				TeamID: 1,
			}},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			repo := createRepo(test.mockFn)

			actual, err := repo.FindByTeamID(context.Background(), test.Param)
			if test.ExpectErr == nil {
				assert.Equal(t, test.Expect, actual)
				assert.Nil(t, err)
			}

			assert.Equal(t, test.ExpectErr, err)
		})
	}
}

func Test_Insert(t *testing.T) {
	testCases := []struct {
		Name      string
		Param     model.PlayerModel
		mockFn    mockFn
		Expect    model.PlayerModel
		ExpectErr error
	}{
		{
			Name: "when_successful",
			Param: model.PlayerModel{
				Name:   "some-player-name",
				TeamID: 1,
			},
			mockFn: func(db sqlmock.Sqlmock) {
				db.ExpectExec(regexp.QuoteMeta("INSERT INTO player (name, team_id) VALUES ($1, $2)")).
					WithArgs("some-player-name", int64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			repo := createRepo(test.mockFn)

			err := repo.Insert(context.Background(), test.Param)

			assert.Equal(t, test.ExpectErr, err)
		})
	}
}
