package store

import "github.com/zhainar/awesomeProject/internal/app/model"

type UserRepository interface {
	Create(u *model.User) error
	FindByEmail(email string) (*model.User, error)
}
