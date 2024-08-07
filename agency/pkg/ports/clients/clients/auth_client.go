package clients

import "agency/internal/user"

type IAuthClient interface {
	GetUserByToken(string) (*user.User, error)
}
