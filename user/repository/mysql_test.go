package mysql_test

import (
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/andhikagama/lmnlo/models/entity"
	"github.com/andhikagama/lmnlo/models/filter"
	userRepo "github.com/andhikagama/lmnlo/user/repository"
)

var mockUser = entity.User{
	ID:       1,
	Email:    `andhika.gama@outlook.com`,
	Password: `aiueo`,
	Address:  `Menteng`,
}

var mockUsers = []*entity.User{
	&mockUser,
}

func TestStore(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`INSERT INTO user`).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := userRepo.NewUserRepository(db)
		err := repo.Store(&mockUser)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-begin", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(fmt.Errorf("Some error"))
		repo := userRepo.NewUserRepository(db)
		err := repo.Store(&mockUser)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-prepare", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`INSERT INTO user`).WillReturnError(fmt.Errorf("Some error"))

		repo := userRepo.NewUserRepository(db)
		err := repo.Store(&mockUser)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-exec", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`INSERT INTO user`).ExpectExec().WillReturnError(fmt.Errorf("Some error"))

		repo := userRepo.NewUserRepository(db)
		err := repo.Store(&mockUser)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-id", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`INSERT INTO user`).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("Some error")))
		mock.ExpectRollback()

		repo := userRepo.NewUserRepository(db)
		err := repo.Store(&mockUser)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			`id`, `email`, `address`,
		}).AddRow(
			mockUsers[0].ID, mockUsers[0].Email, mockUsers[0].Address,
		)

		mock.ExpectQuery(`SELECT (.+) FROM user`).WillReturnRows(rows)

		f := new(filter.User)
		repo := userRepo.NewUserRepository(db)
		res, err := repo.Fetch(f)

		assert.NoError(t, err)
		assert.Equal(t, mockUsers[0].ID, res[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success-with-params", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			`id`, `email`, `address`,
		}).AddRow(
			mockUsers[0].ID, mockUsers[0].Email, mockUsers[0].Address,
		)

		mock.ExpectQuery(`SELECT (.+) FROM user`).WillReturnRows(rows)

		f := new(filter.User)
		f.Email = `andhika.gama@outlook.com`

		repo := userRepo.NewUserRepository(db)
		res, err := repo.Fetch(f)

		assert.NoError(t, err)
		assert.Equal(t, mockUsers[0].ID, res[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success-no-data", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			`id`, `email`, `address`,
		})

		mock.ExpectQuery(`SELECT (.+) FROM user`).WillReturnRows(rows)

		f := new(filter.User)

		repo := userRepo.NewUserRepository(db)
		res, err := repo.Fetch(f)

		assert.NoError(t, err)
		assert.Len(t, res, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			`id`, `email`, `address`,
		}).AddRow(
			mockUsers[0].ID, mockUsers[0].Email, nil,
		).RowError(
			1, fmt.Errorf("row error"),
		)

		mock.ExpectQuery(`SELECT (.+) FROM user`).WillReturnRows(rows)

		f := new(filter.User)

		repo := userRepo.NewUserRepository(db)
		res, err := repo.Fetch(f)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Update(&mockUser)

		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success-no-data", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Update(&mockUser)

		assert.NoError(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-begin", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(fmt.Errorf("Some error"))
		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Update(&mockUser)

		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-prepare", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).WillReturnError(fmt.Errorf("Some error"))

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Update(&mockUser)

		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-exec", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).ExpectExec().WillReturnError(fmt.Errorf("Some error"))

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Update(&mockUser)

		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-affected", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("Some error")))
		mock.ExpectRollback()

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Update(&mockUser)

		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			`id`, `email`, `address`,
		}).AddRow(
			mockUsers[0].ID, mockUsers[0].Email, mockUsers[0].Address,
		)

		mock.ExpectQuery(`SELECT (.+) FROM user`).WillReturnRows(rows)

		repo := userRepo.NewUserRepository(db)
		res, err := repo.GetByID(mockUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, mockUsers[0].ID, res.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success-no-data", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			`id`, `email`, `address`,
		})

		mock.ExpectQuery(`SELECT (.+) FROM user`).WillReturnRows(rows)

		repo := userRepo.NewUserRepository(db)
		res, err := repo.GetByID(mockUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, res.ID, int64(0))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			`id`, `email`, `address`,
		}).AddRow(
			mockUsers[0].ID, mockUsers[0].Email, nil,
		).RowError(
			1, fmt.Errorf("row error"),
		)

		mock.ExpectQuery(`SELECT (.+) FROM user`).WillReturnRows(rows)

		repo := userRepo.NewUserRepository(db)
		res, err := repo.GetByID(mockUser.ID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Delete(mockUser.ID)

		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success-no-data", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Delete(mockUser.ID)

		assert.NoError(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-begin", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(fmt.Errorf("Some error"))
		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Delete(mockUser.ID)

		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-prepare", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).WillReturnError(fmt.Errorf("Some error"))

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Delete(mockUser.ID)

		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-exec", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).ExpectExec().WillReturnError(fmt.Errorf("Some error"))

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Delete(mockUser.ID)

		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error-affected", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(`UPDATE user`).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("Some error")))
		mock.ExpectRollback()

		repo := userRepo.NewUserRepository(db)
		ok, err := repo.Delete(mockUser.ID)

		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
