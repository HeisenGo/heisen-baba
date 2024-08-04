package presenter

import (
	"authservice/internal/user"
	"authservice/protobufs"
)

func UserRegisterToUserDomain(u *protobufs.RegisterRequest) *user.User {
	return &user.User{
		Email:    u.Email,
		Password: u.Password,
	}
}
