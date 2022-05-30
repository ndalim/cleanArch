package delivery

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"net/http/httptest"
	"redditapp/pkg/auth/usecase"
	"redditapp/pkg/user"
	"strings"
	"testing"
)

type TestIn struct {
	params   []interface{}
	needMock bool
	method   string
	url      string
	body     string
}

type TestOut struct {
	code int
}

func TestSignIn(t *testing.T) {

	unitTests := []struct {
		num  int
		name string
		in   TestIn
		out  struct {
			code int
		}
	}{
		{
			num:  1,
			name: "good way",
			in: TestIn{
				params: []interface{}{
					"test_user",
					"12341234",
					&usecase.AuthResult{&user.User{Id: 1, Name: "test_user", HashPass: []byte("pass")}, "token"},
					nil,
				},
				needMock: true,
				method:   "POST",
				url:      "http://localhost:8086/api/login",
				body:     `{"username": "test_user", "password": "12341234"}`,
			},
			out: TestOut{
				code: 200,
			},
		},
		{
			num:  2,
			name: "wrong response body json",
			in: TestIn{
				params:   []interface{}{},
				needMock: false,
				method:   "POST",
				url:      "http://localhost:8086/api/login",
				body:     `{username: "test_user", password: "12341234"}`,
			},
			out: TestOut{
				code: 500,
			},
		},
		{
			num:  3,
			name: "wrong request body",
			in: TestIn{
				params: []interface{}{
					"test_user",
					"12341234",
					nil,
					errors.New("wrong request"),
				},
				needMock: true,
				method:   "POST",
				url:      "http://localhost:8086/api/login",
				body:     `{"username": "test_user", "password": "12341234"}`,
			},
			out: TestOut{
				code: 500,
			},
		},
	}

	for _, tunit := range unitTests {
		ctrlAuth := gomock.NewController(t)
		au := usecase.NewMockAuthInterface(ctrlAuth)

		if tunit.in.needMock {
			if tunit.in.params[3] == nil {
				au.
					EXPECT().
					SignIn(tunit.in.params[0], tunit.in.params[1]).
					Return(tunit.in.params[2].(*usecase.AuthResult), nil)
			} else {
				au.
					EXPECT().
					SignIn(tunit.in.params[0], tunit.in.params[1]).
					Return(nil, tunit.in.params[3].(error))
			}
		}

		handlers := NewAuthHandler(au, nil)

		body := strings.NewReader(tunit.in.body)
		r := httptest.NewRequest(tunit.in.method, tunit.in.url, body)
		w := httptest.NewRecorder()

		handlers.SignIn(w, r)

		res := w.Result()
		if res.StatusCode != tunit.out.code {
			b, _ := io.ReadAll(res.Body)
			t.Errorf("%v %v - want:%v got:%v body:%v",
				tunit.num, tunit.name, tunit.out.code, res.StatusCode, string(b))
		} else {
			fmt.Printf("%v - OK\n", tunit.num)
		}

		ctrlAuth.Finish()
	}
}

func TestSignUp(t *testing.T) {
	unitTests := []struct {
		num  int
		name string
		in   TestIn
		out  struct {
			code int
		}
	}{
		{
			num:  1,
			name: "good way",
			in: TestIn{
				params: []interface{}{
					"test_user",
					"12341234",
					&usecase.AuthResult{&user.User{Id: 1, Name: "test_user", HashPass: []byte("pass")}, "token"},
					nil,
				},
				needMock: true,
				method:   "POST",
				url:      "http://localhost:8086/api/login",
				body:     `{"username": "test_user", "password": "12341234"}`,
			},
			out: TestOut{
				code: 200,
			},
		},
		{
			num:  2,
			name: "wrong response body json",
			in: TestIn{
				params:   []interface{}{},
				needMock: false,
				method:   "POST",
				url:      "http://localhost:8086/api/login",
				body:     `{username: "test_user", password: "12341234"}`,
			},
			out: TestOut{
				code: 500,
			},
		},
		{
			num:  3,
			name: "wrong request body",
			in: TestIn{
				params: []interface{}{
					"test_user",
					"12341234",
					nil,
					errors.New("wrong request"),
				},
				needMock: true,
				method:   "POST",
				url:      "http://localhost:8086/api/login",
				body:     `{"username": "test_user", "password": "12341234"}`,
			},
			out: TestOut{
				code: 500,
			},
		},
	}

	for _, tunit := range unitTests {
		ctrlAuth := gomock.NewController(t)
		au := usecase.NewMockAuthInterface(ctrlAuth)

		if tunit.in.needMock {
			if tunit.in.params[3] == nil {
				au.
					EXPECT().
					SignUp(tunit.in.params[0], tunit.in.params[1]).
					Return(tunit.in.params[2].(*usecase.AuthResult), nil)
			} else {
				au.
					EXPECT().
					SignUp(tunit.in.params[0], tunit.in.params[1]).
					Return(nil, tunit.in.params[3].(error))
			}
		}

		handlers := NewAuthHandler(au, nil)

		body := strings.NewReader(tunit.in.body)
		r := httptest.NewRequest(tunit.in.method, tunit.in.url, body)
		w := httptest.NewRecorder()

		handlers.SignUp(w, r)

		res := w.Result()
		if res.StatusCode != tunit.out.code {
			b, _ := io.ReadAll(res.Body)
			t.Errorf("%v %v - want:%v got:%v body:%v",
				tunit.num, tunit.name, tunit.out.code, res.StatusCode, string(b))
		} else {
			fmt.Printf("%v - OK\n", tunit.num)
		}

		ctrlAuth.Finish()
	}
}
