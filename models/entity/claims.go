package entity

import (
	"github.com/dgrijalva/jwt-go"
)

// Claims ...
type Claims struct {
	User *User `json:"user"`
	jwt.StandardClaims
}
