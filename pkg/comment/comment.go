package comment

import (
	"redditapp/pkg/user"
)

type Comment struct {
	Created string     `json:"created"`
	Author  *user.User `json:"author"`
	Body    string     `json:"body"`
	Id      string     `json:"id"`
	PostId  string     `json:"-"`
}
