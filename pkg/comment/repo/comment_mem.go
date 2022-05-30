package repo

import (
	"errors"
	"redditapp/pkg/comment"
	"strconv"
	"sync"
)

type CommentRepoMem struct {
	List  []comment.Comment
	MaxId int
	Mtx   *sync.RWMutex
}

func NewCommentRepoMem() *CommentRepoMem {
	return &CommentRepoMem{
		List:  make([]comment.Comment, 0, 10),
		Mtx:   &sync.RWMutex{},
		MaxId: 0,
	}
}

func (r *CommentRepoMem) GetByPost(id string) ([]*comment.Comment, error) {
	res := make([]*comment.Comment, 0, 5)

	r.Mtx.RLock()
	for _, v := range r.List {
		if v.PostId != id {
			continue
		}
		res = append(res, &v)
	}
	r.Mtx.RUnlock()

	return res, nil
}

func (r *CommentRepoMem) AddComment(comm comment.Comment) (string, error) {
	r.Mtx.Lock()

	r.MaxId += 1
	comm.Id = strconv.Itoa(r.MaxId)
	r.List = append(r.List, comm)

	r.Mtx.Unlock()
	return comm.Id, nil
}

func (r *CommentRepoMem) Delete(postId, id, user string) error {
	r.Mtx.Lock()
	defer r.Mtx.Unlock()

	for i, v := range r.List {
		if v.Id != id {
			continue
		}
		if v.PostId != postId {
			continue
		}
		if v.Author.Name != user {
			return errors.New("current user is not an author of the comment")
		}
		r.List = append(r.List[:i], r.List[i+1:]...)
		return nil
	}
	return errors.New("comment has not found")
}
