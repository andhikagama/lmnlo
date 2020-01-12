package user

import (
	"github.com/andhikagama/lmnlo/models/entity"
)

// Repository represents database manipulation
type Repository interface {
	Store(*entity.User) error
}

// Usecase represents business logic
type Usecase interface {
	Register(g *entity.User) error
}
