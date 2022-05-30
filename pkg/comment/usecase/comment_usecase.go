package usecase

import (
	"redditapp/pkg/comment"
	"redditapp/pkg/comment/repo"
	"redditapp/pkg/user"
	"time"
)

type CommentUsecase struct {
	CommentRepo repo.CommentRepoInterface
}

func NewCommentUsecase(repo repo.CommentRepoInterface) *CommentUsecase {
	return &CommentUsecase{
		CommentRepo: repo,
	}
}
func (u *CommentUsecase) GetByPost(id string) ([]*comment.Comment, error) {
	return u.CommentRepo.GetByPost(id)
}

func (u *CommentUsecase) AddComment(postId string, body string, author *user.User) (string, error) {

	comm := &comment.Comment{
		PostId:  postId,
		Body:    body,
		Author:  author,
		Created: time.Now().Format(time.RFC3339),
	}
	id, err := u.CommentRepo.AddComment(*comm)
	return id, err
}

func (u *CommentUsecase) Delete(postId string, id string, user string) error {
	return u.CommentRepo.Delete(postId, id, user)
}
