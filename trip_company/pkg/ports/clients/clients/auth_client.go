package clients

import "tripcompanyservice/internal/user"


type IAuthClient interface {
	GetUserByToken(string) (*user.User, error)
}
