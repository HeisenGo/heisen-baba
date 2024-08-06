package clients

import "hotel/internal/user"

type IAuthClient interface {
	GetUserByToken(string) (*user.User, error)
}
