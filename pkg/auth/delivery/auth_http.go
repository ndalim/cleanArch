package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"redditapp/pkg/auth"
	"redditapp/pkg/user/repo"
)

const currUserKey = "currUser"

type jwtResponse struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	UserRepo    repo.UserRepoInterface
	AuthUsecase auth.AuthInterface
}

func NewAuthHandler(a auth.AuthInterface, r repo.UserRepoInterface) *AuthHandler {
	return &AuthHandler{
		UserRepo:    r,
		AuthUsecase: a,
	}
}

type authData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func bodyRequest(r *http.Request, ad *authData) error {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return err
	}

	return json.Unmarshal(body, ad)
}

func returnSession(token string, w http.ResponseWriter, ad *authData) {
	res, err := json.Marshal(&jwtResponse{token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Write(res)
}

func (a *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {

	ad := &authData{}
	err := bodyRequest(r, ad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sign, err := a.AuthUsecase.SignIn(ad.Username, ad.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnSession(sign.Token, w, ad)
	return
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	ad := &authData{}
	err := bodyRequest(r, ad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sign, err := a.AuthUsecase.SignUp(ad.Username, ad.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnSession(sign.Token, w, ad)
	return
}

//test

func (a *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	//vars := mux.Vars(r)
	//usr, err := a.UserRepo.GetUser(vars["name"])
	//
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//res, err := json.Marshal(usr)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//w.WriteHeader(http.StatusOK)
	//w.Write(res)
}

func (a *AuthHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	//usrs, err := a.UserRepo.GetAll()
	//
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Write([]byte(err.Error()))
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//res, err := json.Marshal(usrs)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Write([]byte(err.Error()))
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//w.WriteHeader(http.StatusOK)
	//w.Write(res)
}
