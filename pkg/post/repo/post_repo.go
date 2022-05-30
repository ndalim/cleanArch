package repo

import (
	"redditapp/pkg/post"
	"redditapp/pkg/user"
)

type PostRepoInterface interface {
	GetAll() ([]post.Post, error)
	GetByCategory(c post.PCategory) ([]post.Post, error)
	GetByUser(user string) ([]post.Post, error)
	GetPost(id string) (*post.Post, error)
	ShowPost(id string) (*post.Post, error)
	AddPost(p *post.Post, author *user.User) (id string, err error)
	UpVote(id string, user string) error
	DownVote(id string, user string) error
	UnVote(id string, user string) error
	Delete(id string, user string) error
}
