package user

import (
	"github.com/andhikagama/lmnlo/models/entity"
	"github.com/andhikagama/lmnlo/models/filter"
)

// Repository represents database manipulation
type Repository interface {
	Store(usr *entity.User) error
	Fetch(f *filter.User) ([]*entity.User, error)
	Update(usr *entity.User) (bool, error)
	GetByID(id int64) (*entity.User, error)
	Delete(id int64) (bool, error)
	InsertToken(uid int64, token string) error
	ValidateToken(token string) (bool, error)
}

// Usecase represents business logic
type Usecase interface {
	Register(usr *entity.User) error
	Fetch(f *filter.User) ([]*entity.User, error)
	Update(usr *entity.User) error
	GetByID(id int64) (*entity.User, error)
	Delete(id int64) error
	PartialUpdate(id int64, byteFacility []byte) (*entity.User, error)
	Login(u *entity.User) (*entity.User, error)
}
