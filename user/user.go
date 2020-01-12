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
}

// Usecase represents business logic
type Usecase interface {
	Register(usr *entity.User) error
	Fetch(f *filter.User) ([]*entity.User, error)
	Update(usr *entity.User) error
}
