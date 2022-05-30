package repo

import (
	"errors"
	"redditapp/pkg/user"
)

var ErrorGettingUser = errors.New("user not found")

type UserRepoInterface interface {
	GetUser(name string) (*user.User, error)
	GetAll() ([]*user.User, error)
	Create(string, []byte) (*user.User, error)
}
