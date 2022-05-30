package usecase

import (
	"bytes"
	"crypto/sha1"
	"errors"
	ar "redditapp/pkg/auth/repo"
	"redditapp/pkg/user"
	ur "redditapp/pkg/user/repo"
)

var (
	ErrorWrongSingIn = errors.New("bad login/password")
	ErrorSession     = errors.New("wrong session")
)

type AuthUsecase struct {
	UserRepo ur.UserRepoInterface
	AuthRepo ar.AuthRepoInterface
	Salt     string
}

type AuthResult struct {
	User  *user.User
	Token string
}

func authResult(usr *user.User, token string) *AuthResult {
	return &AuthResult{
		User:  usr,
		Token: token,
	}
}

func NewAuthUsecase(repo ur.UserRepoInterface, authRepo ar.AuthRepoInterface, salt string) *AuthUsecase {
	a := &AuthUsecase{
		UserRepo: repo,
		AuthRepo: authRepo,
		Salt:     salt,
	}
	return a
}

func (a *AuthUsecase) CheckPass(plainPass string, hashCheck []byte) bool {
	hash := makeHashPass(plainPass, a.Salt)
	return bytes.Equal(hash, hashCheck)
}

func (a *AuthUsecase) SignIn(username, password string) (*AuthResult, error) {

	usr, err := a.UserRepo.GetUser(username)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Usecase SignIn 1: ", usr, err)

	if !a.CheckPass(password, usr.HashPass) {
		return nil, ErrorWrongSingIn
	}
	//fmt.Println("Usecase SignIn 2: ", password, usr.HashPass)

	token, err := a.AuthRepo.Session(usr)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Usecase SignIn 3: ", token, err)

	return authResult(usr, token), nil
}

func (a *AuthUsecase) SignUp(username, password string) (*AuthResult, error) {

	hashPass := makeHashPass(password, a.Salt)
	usr, err := a.UserRepo.Create(username, hashPass)
	if err != nil {
		return nil, err
	}

	token, err := a.AuthRepo.Session(usr)
	if err != nil {
		return nil, err
	}

	return authResult(usr, token), nil
}

func (a *AuthUsecase) Check(tokenString string) (*user.User, error) {

	name, err := a.AuthRepo.Check(tokenString)

	if err != nil {
		return nil, err
	}

	usr, err := a.UserRepo.GetUser(name)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func makeHashPass(password, salt string) []byte {
	pass := sha1.New()
	pass.Write([]byte(password))
	pass.Write([]byte(salt))

	return pass.Sum(nil)
}
