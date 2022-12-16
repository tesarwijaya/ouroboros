package repository_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tesarwijaya/ouroboros/internal/domain/team/model"
	"github.com/tesarwijaya/ouroboros/internal/domain/team/repository"
)

type mockFn func(db sqlmock.Sqlmock)

func createRepo(mockFn mockFn) repository.TeamRepository {
	db, mock, _ := sqlmock.New()

	mockFn(mock)
	repo := repository.NewTeamReposity(repository.TeamRepositoryImpl{
		Db: db,
	})

	return repo
}

func Test_FindAll(t *testing.T) {
	testCases := []struct {
		Name        string
		MockFn      mockFn
		Expected    []model.TeamModel
		ExpectedErr string
	}{
		{
			Name: "when_data_present",
			MockFn: func(db sqlmock.Sqlmock) {
				db.ExpectQuery(regexp.QuoteMeta("SELECT * FROM team")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(int64(1), "some-team-name"),
					)
			},
			Expected: []model.TeamModel{{
				ID:   1,
				Name: "some-team-name",
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
		Name        string
		Param       int64
		MockFn      mockFn
		Expected    model.TeamModel
		ExpectedErr string
	}{
		{
			Name:  "when_data_present",
			Param: 1,
			MockFn: func(db sqlmock.Sqlmock) {
				db.ExpectQuery(regexp.QuoteMeta("SELECT * FROM team WHERE id = $1")).WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(int64(1), "some-team-name"),
					)
			},
			Expected: model.TeamModel{
				ID:   1,
				Name: "some-team-name",
			},
		},
	}

	for _, test := range testCases {
		repo := createRepo(test.MockFn)

		actual, err := repo.FindByID(context.Background(), test.Param)
		if test.ExpectedErr != "" {
			assert.EqualError(t, err, test.ExpectedErr)
		} else {
			assert.Equal(t, test.Expected, actual)
			assert.Nil(t, err)
		}
	}
}

func Test_Insert(t *testing.T) {
	testCases := []struct {
		Name        string
		Param       model.TeamModel
		MockFn      mockFn
		ExpectedErr string
	}{
		{
			Name: "when_successful",
			Param: model.TeamModel{
				Name: "some-team-name",
			},
			MockFn: func(db sqlmock.Sqlmock) {
				db.ExpectExec(regexp.QuoteMeta("INSERT INTO team (name) VALUES ($1)")).WithArgs("some-team-name").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, test := range testCases {
		repo := createRepo(test.MockFn)

		err := repo.Insert(context.Background(), test.Param)

		if test.ExpectedErr != "" {
			assert.EqualError(t, err, test.ExpectedErr)
		} else {
			assert.Nil(t, err)
		}

	}
}
