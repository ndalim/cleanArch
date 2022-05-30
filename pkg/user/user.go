package user

import "errors"

var (
	ErrCreateUser = errors.New("can't create user")
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"username"`
	HashPass []byte `json:"-"`
}
