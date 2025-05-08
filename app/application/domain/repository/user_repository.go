package repository

import "github.com/softorwhite/lambda-o11y/app/application/domain/model"

type UserRepository interface {
	GetUser(userID string) (*model.User, error)
}
