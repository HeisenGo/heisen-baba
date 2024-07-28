package clients

import "terminalpathservice/internal/user"

type IAuthClient interface {
	GetUserByToken(string) (*user.User, error)
}
