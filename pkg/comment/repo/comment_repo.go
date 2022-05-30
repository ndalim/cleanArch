package repo

import (
	"redditapp/pkg/comment"
)

type CommentRepoInterface interface {
	GetByPost(postId string) ([]*comment.Comment, error)
	AddComment(comm comment.Comment) (id string, err error)
	Delete(postId, id, user string) error
}
