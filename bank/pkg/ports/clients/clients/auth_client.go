package clients

import "bankservice/internal/user"

type IAuthClient interface {
	GetUserByToken(string) (*user.User, error)
}
