package usecase

import (
	"redditapp/pkg/post"
	pr "redditapp/pkg/post/repo"
)

type PostUsecase struct {
	PostRepo pr.PostRepoInterface
}

func NewPostUsecase(repo pr.PostRepoInterface) *PostUsecase {
	p := &PostUsecase{
		PostRepo: repo,
	}
	return p
}

func (uc *PostUsecase) GetAll() ([]post.Post, error) {
	return uc.PostRepo.GetAll()
}

func (uc *PostUsecase) GetPost(id string) (*post.Post, error) {
	return uc.PostRepo.GetPost(id)
}

func (uc *PostUsecase) ShowPost(id string) (*post.Post, error) {
	return uc.PostRepo.ShowPost(id)
}

func (uc *PostUsecase) GetByCategory(categoryName string) ([]post.Post, error) {
	category := post.PCategory(categoryName)
	return uc.PostRepo.GetByCategory(category)
}

func (uc *PostUsecase) GetByUser(user string) ([]post.Post, error) {
	return uc.PostRepo.GetByUser(user)
}

func (uc *PostUsecase) UpVote(id string, user string) error {
	return uc.PostRepo.UpVote(id, user)
}

func (uc *PostUsecase) DownVote(id string, user string) error {
	return uc.PostRepo.DownVote(id, user)
}

func (uc *PostUsecase) UnVote(id string, user string) error {
	return uc.PostRepo.UnVote(id, user)
}

func (uc *PostUsecase) Delete(id string, usr string) error {
	return uc.PostRepo.Delete(id, usr)
}
