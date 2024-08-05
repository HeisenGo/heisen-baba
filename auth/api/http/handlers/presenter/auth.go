package presenter

import (
	"authservice/internal/user"
)

type UserRegisterReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func UserRegisterToUserDomain(up *UserRegisterReq) *user.User {
	return &user.User{
		Email:    up.Email,
		Password: up.Password,
	}
}
