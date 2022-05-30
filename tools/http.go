package tools

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditapp/pkg/user"
)

type sessionKey string

var CurrentUserKey sessionKey = "currentUser"

func CurrentUser(r *http.Request) (*user.User, bool) {

	val := r.Context().Value(CurrentUserKey)
	usr, ok := val.(*user.User)

	return usr, ok
}

func RequestBody(r *http.Request, value interface{}) error {
	body := r.Body
	defer r.Body.Close()

	bodyBite, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBite, value)
}
