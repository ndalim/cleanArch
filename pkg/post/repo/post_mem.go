package repo

import (
	"errors"
	"redditapp/pkg/post"
	"redditapp/pkg/user"
	"redditapp/tools"
	"strconv"
	"sync"
	"time"
)

type PostRepoMem struct {
	List    []post.Post
	IndexId map[string]*post.Post
	MaxId   int
	Mtx     *sync.RWMutex
}

func NewPostRepoMem() *PostRepoMem {
	p := &PostRepoMem{
		List:    make([]post.Post, 0, 10),
		MaxId:   0,
		Mtx:     &sync.RWMutex{},
		IndexId: make(map[string]*post.Post, 10),
	}

	return p
}

func (r *PostRepoMem) GetAll() ([]post.Post, error) {

	res := make([]post.Post, len(r.List))

	r.Mtx.RLock()

	for i, v := range r.List {
		res[i] = v
		i++
	}
	r.Mtx.RUnlock()

	return res, nil
}

func (r *PostRepoMem) GetByCategory(c post.PCategory) ([]post.Post, error) {

	res := make([]post.Post, 0, 5)

	r.Mtx.RLock()
	for _, v := range r.List {
		if v.Category != c {
			continue
		}
		res = append(res, v)
	}
	r.Mtx.RUnlock()

	return res, nil
}

func (r *PostRepoMem) GetByUser(user string) ([]post.Post, error) {

	res := make([]post.Post, 0, 5)

	r.Mtx.RLock()
	for _, v := range r.List {
		if v.Author.Name != user {
			continue
		}
		res = append(res, v)
	}
	r.Mtx.RUnlock()

	return res, nil
}

func (r *PostRepoMem) GetPost(id string) (*post.Post, error) {

	r.Mtx.RLock()
	defer r.Mtx.RUnlock()

	pst, ok := r.IndexId[id]
	if !ok {
		return nil, post.ErrorPostNotFound
	}

	return pst, nil
}

func (r *PostRepoMem) ShowPost(id string) (*post.Post, error) {

	r.Mtx.Lock()
	defer r.Mtx.Unlock()

	pst, ok := r.IndexId[id]
	if !ok {
		return nil, post.ErrorPostNotFound
	}
	pst.Views += 1
	return pst, nil
}

func (r *PostRepoMem) AddPost(pst *post.Post, author *user.User) (string, error) {

	r.Mtx.Lock()

	r.MaxId += 1

	pst.Id = strconv.Itoa(r.MaxId)
	pst.Created = time.Now().Format(time.RFC3339)
	pst.Votes = make([]post.Vote, 0)
	pst.Author = author

	r.List = append(r.List, *pst)
	r.IndexId[pst.Id] = pst

	r.Mtx.Unlock()

	return pst.Id, nil
}

func (r *PostRepoMem) SetVote(id string, user string, vote int) error {

	pst, err := r.GetPost(id)
	if err != nil {
		return err
	}

	r.Mtx.Lock()
	defer r.Mtx.Unlock()
	defer tools.CalculateUpvotePercentage(pst)

	for i, v := range pst.Votes {
		if v.User != user {
			continue
		}
		pst.Votes[i].Vote = vote
		return nil
	}

	pst.Votes = append(pst.Votes, post.Vote{
		User: user,
		Vote: vote,
	})
	return nil
}

func (r *PostRepoMem) UpVote(id string, user string) error {
	return r.SetVote(id, user, 1)
}

func (r *PostRepoMem) DownVote(id string, user string) error {
	return r.SetVote(id, user, -1)
}

func (r *PostRepoMem) UnVote(id string, user string) error {
	pst, err := r.GetPost(id)
	if err != nil {
		return err
	}

	r.Mtx.Lock()
	defer r.Mtx.Unlock()
	defer tools.CalculateUpvotePercentage(pst)

	for i, v := range pst.Votes {
		if v.User != user {
			continue
		}
		pst.Votes = append(pst.Votes[:i], pst.Votes[i+1:]...)
		return nil
	}
	return nil

}

func (r *PostRepoMem) Delete(id string, user string) error {
	r.Mtx.Lock()
	defer r.Mtx.Unlock()

	for i, v := range r.List {
		if v.Id != id {
			continue
		}
		if v.Author.Name != user {
			return errors.New("current user is not an author of the post")
		}

		r.List = append(r.List[:i], r.List[i+1:]...)
		return nil
	}
	return errors.New("post has not found")
}
