package mappers

import (
	"github.com/google/uuid"
	"terminalpathservice/internal/user"
	"terminalpathservice/protobufs"
)

func UserClaimsToDomain(p *protobufs.GetUserByTokenResponse) (*user.User, error) {
	u, err := uuid.Parse(p.UserId)
	if err != nil {
		return nil, err
	}
	return &user.User{
		ID:      u,
		IsAdmin: p.IsAdmin,
	}, nil
}
