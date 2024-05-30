package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

type (
	AuthRequestDTO struct {
		UserName string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

type CustomValidator struct {
	Validator *validator.Validate
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type JwtCustomClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}
