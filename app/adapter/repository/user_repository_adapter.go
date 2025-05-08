package repository

import (
	"github.com/softorwhite/lambda-o11y/app/application/domain/model"
	"github.com/softorwhite/lambda-o11y/app/application/domain/repository"
)

type UserRepositoryAdapter struct{}

func NewUserRepositoryAdapter() repository.UserRepository {
	return &UserRepositoryAdapter{}
}
func (u *UserRepositoryAdapter) GetUser(userID string) (*model.User, error) {
	return &model.User{
		ID:   userID,
		Name: "John Doe",
	}, nil
}
