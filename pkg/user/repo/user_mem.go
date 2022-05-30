package repo

import (
	"redditapp/pkg/user"
	"sync"
)

type UserRepoMem struct {
	List  map[string]*user.User
	Mtx   *sync.RWMutex
	MaxId int
}

func NewUserRepoMem() UserRepoInterface {
	return &UserRepoMem{
		List:  map[string]*user.User{},
		Mtx:   &sync.RWMutex{},
		MaxId: 0,
	}
}

func (u *UserRepoMem) GetUser(name string) (*user.User, error) {

	u.Mtx.RLock()
	defer u.Mtx.RUnlock()

	search, ok := u.List[name]
	if !ok {
		return nil, ErrorGettingUser
	}
	return search, nil
}

func (u *UserRepoMem) GetAll() ([]*user.User, error) {

	res := make([]*user.User, len(u.List))

	u.Mtx.RLock()
	i := 0
	for _, v := range u.List {
		res[i] = v
		i++
	}
	u.Mtx.RUnlock()

	return res, nil
}

func (u *UserRepoMem) Create(name string, pass []byte) (*user.User, error) {

	if _, err := u.GetUser(name); err == nil {
		return nil, user.ErrCreateUser
	}

	usr := &user.User{
		Name: name,
	}

	u.Mtx.Lock()

	u.MaxId += 1
	u.List[name] = usr

	usr.Id = u.MaxId
	usr.HashPass = pass

	u.Mtx.Unlock()

	return usr, nil
}
