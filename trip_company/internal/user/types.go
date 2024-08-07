package user

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type User struct {
	ID      uuid.UUID
	IsAdmin bool
}
