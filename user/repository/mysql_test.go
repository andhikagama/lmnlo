package mysql_test

import (
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/andhikagama/lmnlo/models/entity"
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
