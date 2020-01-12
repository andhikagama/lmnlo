package usecase_test

import (
	"errors"
	"testing"

	"github.com/andhikagama/lmnlo/models/filter"
	"github.com/andhikagama/lmnlo/models/response"

	"github.com/andhikagama/lmnlo/models/entity"
	"github.com/andhikagama/lmnlo/user/mocks"
	"github.com/andhikagama/lmnlo/user/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	mockUserRepo := new(mocks.Repository)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Fetch", mock.AnythingOfType("*filter.User")).Return(make([]*entity.User, 0), nil).Once()
		mockUserRepo.On("Store", mock.AnythingOfType("*entity.User")).Return(nil).Once()
		u := usecase.NewUserUsecase(mockUserRepo)

		err := u.Register(&mockUser)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("already-exist", func(t *testing.T) {
		mockUserRepo.On("Fetch", mock.AnythingOfType("*filter.User")).Return(mockUsers, nil).Once()
		u := usecase.NewUserUsecase(mockUserRepo)

		err := u.Register(&mockUser)

		assert.Error(t, err)
		assert.EqualError(t, err, response.ErrAlreadyExist.Error())
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockUserRepo.On("Fetch", mock.AnythingOfType("*filter.User")).Return(make([]*entity.User, 0), nil).Once()
		mockUserRepo.On("Store", mock.AnythingOfType("*entity.User")).Return(errors.New(`error`)).Once()
		u := usecase.NewUserUsecase(mockUserRepo)

		err := u.Register(&mockUser)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestFetch(t *testing.T) {
	mockUserRepo := new(mocks.Repository)

	t.Run("success", func(t *testing.T) {
		f := new(filter.User)
		mockUserRepo.On("Fetch", mock.AnythingOfType("*filter.User")).Return(mockUsers, nil).Once()
		u := usecase.NewUserUsecase(mockUserRepo)

		res, err := u.Fetch(f)

		assert.NoError(t, err)
		assert.Equal(t, mockUsers, res)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("success-no-data", func(t *testing.T) {
		f := new(filter.User)
		mockEmptyUsers := make([]*entity.User, 0)
		mockUserRepo.On("Fetch", mock.AnythingOfType("*filter.User")).Return(mockEmptyUsers, nil).Once()
		u := usecase.NewUserUsecase(mockUserRepo)

		res, err := u.Fetch(f)

		assert.NoError(t, err)
		assert.Equal(t, mockEmptyUsers, res)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		f := new(filter.User)

		mockUserRepo.On("Fetch", mock.AnythingOfType("*filter.User")).Return(nil, errors.New(`Error`)).Once()
		u := usecase.NewUserUsecase(mockUserRepo)

		res, err := u.Fetch(f)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockUserRepo := new(mocks.Repository)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(true, nil).Once()
		u := usecase.NewUserUsecase(mockUserRepo)

		err := u.Update(&mockUser)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(false, errors.New(`error`)).Once()
		u := usecase.NewUserUsecase(mockUserRepo)

		err := u.Update(&mockUser)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("no-data", func(t *testing.T) {
		mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(false, nil).Once()
		u := usecase.NewUserUsecase(mockUserRepo)

		err := u.Update(&mockUser)

		assert.Error(t, err)
		assert.Equal(t, response.ErrNotFound, err)
		mockUserRepo.AssertExpectations(t)
	})
}

// func TestDelete(t *testing.T) {
// 	mockUserRepo := new(mocks.Repository)

// 	t.Run("success", func(t *testing.T) {
// 		mockUserRepo.On("Delete", mock.AnythingOfType("int64")).Return(true, nil).Once()
// 		u := usecase.NewUserUsecase(mockUserRepo)

// 		err := u.Delete(mockUser.ID)

// 		assert.NoError(t, err)
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("error", func(t *testing.T) {
// 		mockUserRepo.On("Delete", mock.AnythingOfType("int64")).Return(false, errors.New(`error`)).Once()
// 		u := usecase.NewUserUsecase(mockUserRepo)

// 		err := u.Delete(mockUser.ID)

// 		assert.Error(t, err)
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("no-data", func(t *testing.T) {
// 		mockUserRepo.On("Delete", mock.AnythingOfType("int64")).Return(false, nil).Once()
// 		u := usecase.NewUserUsecase(mockUserRepo)

// 		err := u.Delete(mockUser.ID)

// 		assert.Error(t, err)
// 		assert.Equal(t, response.ErrNotFound, err)
// 		mockUserRepo.AssertExpectations(t)
// 	})
// }
