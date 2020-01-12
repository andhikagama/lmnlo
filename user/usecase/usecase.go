package usecase

import (
	"encoding/json"
	"time"

	patch "gopkg.in/evanphx/json-patch.v4"

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

// PartialUpdate ...
func (u *userUsecase) PartialUpdate(id int64, byteObj []byte) (*entity.User, error) {
	existingUser, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return new(entity.User), response.ErrNotFound
	}

	jsonTarget, _ := json.Marshal(existingUser)

	patchObj, err := patch.DecodePatch(byteObj)
	if err != nil {
		return nil, err
	}

	jsonTarget, err = patchObj.Apply(jsonTarget)
	if err != nil {
		return nil, err
	}

	updatedUser := new(entity.User)
	err = json.Unmarshal(jsonTarget, updatedUser)

	if err != nil {
		return nil, err
	}

	ok, err := u.userRepo.Update(updatedUser)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, response.ErrNotFound
	}

	return updatedUser, nil
}

// Login ...
func (u *userUsecase) Login(usr *entity.User) (*entity.User, error) {
	encryptedPass, _ := helper.EncryptToString(usr.Password)
	usr.Password = encryptedPass

	f := new(filter.User)
	f.Email = usr.Email
	f.Password = usr.Password
	f.Num = 1

	usrs, err := u.userRepo.Fetch(f)
	if len(usrs) == 0 {
		return nil, response.ErrLogin
	}

	if err != nil {
		return nil, err
	}

	usr = usrs[0]
	usr.Password = ``

	cc := new(entity.Claims)
	cc.User = usr
	cc.IssuedAt = time.Now().Unix()
	cc.ExpiresAt = time.Now().AddDate(0, 1, 0).Unix()

	token := helper.GenerateTokenString(cc)
	usr.Token = token

	err = u.userRepo.InsertToken(usr.ID, token)
	if err != nil {
		return nil, err
	}

	ok, err := u.userRepo.Update(usr)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return usr, nil
}
