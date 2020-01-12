package usecase

import (
	"github.com/andhikagama/lmnlo/helper"
	"github.com/andhikagama/lmnlo/models/entity"
	"github.com/andhikagama/lmnlo/models/filter"
	"github.com/andhikagama/lmnlo/models/response"
	"github.com/andhikagama/lmnlo/user"
)

type userUsecase struct {
	userRepo user.Repository
}

// NewUserUsecase ...
func NewUserUsecase(
	r user.Repository,
) user.Usecase {
	return &userUsecase{
		r,
	}
}

// Register ...
func (u *userUsecase) Register(usr *entity.User) error {
	f := new(filter.User)
	f.Email = usr.Email
	f.Num = 1

	usrs, err := u.userRepo.Fetch(f)
	if err != nil {
		return err
	}

	if len(usrs) != 0 {
		return response.ErrAlreadyExist
	}

	encryptedPass, err := helper.EncryptToString(usr.Password)
	if err != nil {
		return err
	}
	usr.Password = encryptedPass

	return u.userRepo.Store(usr)
}

// Fetch ...
func (u *userUsecase) Fetch(f *filter.User) ([]*entity.User, error) {
	return u.userRepo.Fetch(f)
}

// Update ...
func (u *userUsecase) Update(usr *entity.User) error {
	ok, err := u.userRepo.Update(usr)

	if err != nil {
		return err
	}

	if !ok {
		return response.ErrNotFound
	}

	return nil
}

// GetByID ...
func (u *userUsecase) GetByID(id int64) (*entity.User, error) {
	return u.userRepo.GetByID(id)
}

// Delete ...
func (u *userUsecase) Delete(id int64) error {
	ok, err := u.userRepo.Delete(id)

	if err != nil {
		return err
	}

	if !ok {
		return response.ErrNotFound
	}

	return nil
}
