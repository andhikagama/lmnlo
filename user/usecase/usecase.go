package usecase

import (
	"github.com/andhikagama/lmnlo/helper"
	"github.com/andhikagama/lmnlo/models/entity"
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

	encryptedPass, err := helper.EncryptToString(usr.Password)
	if err != nil {
		return err
	}
	usr.Password = encryptedPass

	return u.userRepo.Store(usr)
}
