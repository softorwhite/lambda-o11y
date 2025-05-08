package usecase

import (
	"fmt"

	a_repository "github.com/softorwhite/lambda-o11y/app/adapter/repository"
	"github.com/softorwhite/lambda-o11y/app/application/domain/model"
	"github.com/softorwhite/lambda-o11y/app/application/domain/repository"
)

type UserUseCase struct {
	r func() repository.UserRepository
}

func NewUserUseCase() *UserUseCase {
	return &UserUseCase{
		r: func() repository.UserRepository {
			return a_repository.NewUserRepositoryAdapter()
		},
	}
}

func (u *UserUseCase) GetUser(userID string) (*model.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}
	user, err := u.r().GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
