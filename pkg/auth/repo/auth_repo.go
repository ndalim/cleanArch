package repo

import "redditapp/pkg/user"

type AuthRepoInterface interface {
	Check(string) (string, error)
	Session(user *user.User) (string, error)
}
