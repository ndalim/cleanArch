package auth

import (
	"redditapp/pkg/auth/usecase"
	"redditapp/pkg/user"
)

type AuthModel struct {
	Usecase usecase.AuthUsecase
	Secret  string
	Salt    string
}

type AuthInterface interface {
	SignIn(username, password string) (*usecase.AuthResult, error)
	SignUp(username, password string) (*usecase.AuthResult, error)
	Check(string) (*user.User, error)
}
