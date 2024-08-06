package mappers

import (
	"authservice/internal/user"
	"authservice/pkg/adapters/storage/entities"
)

func UserEntityToDomain(entity *entities.User) *user.User {
	return &user.User{
		ID:       entity.ID,
		Email:    entity.Email,
		Password: entity.Password,
		IsAdmin:  entity.IsAdmin,
	}
}
func userEntityToDomain(entity entities.User) user.User {
	return user.User{
		ID:       entity.ID,
		Email:    entity.Email,
		Password: entity.Password,
	}
}

func UserDomainToEntity(domainUser *user.User) *entities.User {
	return &entities.User{
		Email:    domainUser.Email,
		Password: domainUser.Password,
	}
}
