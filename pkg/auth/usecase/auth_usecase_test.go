package usecase

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	ar "redditapp/pkg/auth/repo"
	"redditapp/pkg/user"
	ur "redditapp/pkg/user/repo"
	"testing"
)

var pass = "pass"
var salt = "salt"

//pass = "pass", salt = "salt"
var goodHash = []byte{17, 245, 99, 159, 34, 82, 81, 85, 203, 11, 67, 87, 62, 228, 33, 40, 56, 199, 141, 135}

func TestSignIn(t *testing.T) {

	usr := &user.User{Id: 1, Name: "alex", HashPass: goodHash}
	result := &AuthResult{usr, "token"}
	err := errors.New("wrong test")

	unitTests := []struct {
		num      int
		name     string
		getUser1 string
		getUser2 *user.User
		getUser3 error
		errPass  bool
		session1 *user.User
		session2 string
		session3 error
		out      *AuthResult
		err      error
		isErr    bool
	}{
		{
			num:      1,
			name:     "good way",
			getUser1: "alex",
			getUser2: usr,
			getUser3: nil,
			errPass:  false,
			session1: usr,
			session2: "token",
			session3: nil,
			out:      result,
			err:      nil,
			isErr:    false,
		},
		{
			num:      2,
			name:     "wrong pass",
			getUser1: "alex",
			getUser2: &user.User{Id: 1, Name: "alex", HashPass: []byte("wrong pass")},
			getUser3: nil,
			errPass:  true,
			session1: usr,
			session2: "",
			session3: err,
			out:      nil,
			err:      err,
			isErr:    true,
		},
		{
			num:      3,
			name:     "getUser error",
			getUser1: "alex",
			getUser2: nil,
			getUser3: err,
			errPass:  true,
			session1: nil,
			session2: "",
			session3: nil,
			out:      nil,
			err:      err,
			isErr:    true,
		},
		{
			num:      4,
			name:     "session error",
			getUser1: "alex",
			getUser2: usr,
			getUser3: nil,
			errPass:  false,
			session1: usr,
			session2: "",
			session3: err,
			out:      nil,
			err:      err,
			isErr:    true,
		},
	}

	for _, test := range unitTests {
		ctrl := gomock.NewController(t)

		uRepo := ur.NewMockUserRepoInterface(ctrl)
		uRepo.
			EXPECT().
			GetUser(test.getUser1).
			Return(test.getUser2, test.getUser3)

		aRepo := ar.NewMockAuthRepoInterface(ctrl)
		if !test.errPass {
			aRepo.
				EXPECT().
				Session(test.session1).
				Return(test.session2, test.session3)
		}

		au := NewAuthUsecase(uRepo, aRepo, salt)
		res, err := au.SignIn("alex", pass)

		if test.isErr {
			if err == nil {
				t.Errorf("%v %v - got:%v, want:%v", test.num, test.name, res, test.out)
			}
		} else {
			if fmt.Sprint(res) != fmt.Sprint(test.out) {
				t.Errorf("%v %v - got:%v, want:%v", test.num, test.name, res, test.out)
			}
		}
		ctrl.Finish()
	}

}

func TestSignUp(t *testing.T) {

	usr := &user.User{Id: 1, Name: "alex", HashPass: goodHash}
	result := &AuthResult{usr, "token"}
	err := errors.New("wrong test")

	unitTests := []struct {
		num      int
		name     string
		create1  string
		create2  []byte
		create3  *user.User
		create4  error
		errPass  bool
		session1 *user.User
		session2 string
		session3 error
		out      *AuthResult
		err      error
		isErr    bool
	}{
		{
			num:      1,
			name:     "good way",
			create1:  "alex",
			create2:  goodHash,
			create3:  usr,
			create4:  nil,
			errPass:  false,
			session1: usr,
			session2: "token",
			session3: nil,
			out:      result,
			err:      nil,
			isErr:    false,
		},
		{
			num:      2,
			name:     "create error",
			create1:  "alex",
			create2:  goodHash,
			create3:  nil,
			create4:  err,
			errPass:  true,
			session1: usr,
			session2: "",
			session3: err,
			out:      nil,
			err:      err,
			isErr:    true,
		},
		{
			num:      3,
			name:     "session error",
			create1:  "alex",
			create2:  goodHash,
			create3:  usr,
			create4:  nil,
			errPass:  false,
			session1: usr,
			session2: "",
			session3: err,
			out:      nil,
			err:      err,
			isErr:    true,
		},
	}

	for _, test := range unitTests {
		ctrl := gomock.NewController(t)

		uRepo := ur.NewMockUserRepoInterface(ctrl)
		uRepo.
			EXPECT().
			Create(test.create1, test.create2).
			Return(test.create3, test.create4)

		aRepo := ar.NewMockAuthRepoInterface(ctrl)
		if !test.errPass {
			aRepo.
				EXPECT().
				Session(test.session1).
				Return(test.session2, test.session3)
		}

		au := NewAuthUsecase(uRepo, aRepo, salt)
		res, err := au.SignUp("alex", pass)

		if test.isErr {
			if err == nil {
				t.Errorf("%v %v - got:%v, want:%v", test.num, test.name, res, test.out)
			}
		} else {
			if fmt.Sprint(res) != fmt.Sprint(test.out) {
				t.Errorf("%v %v - got:%v, want:%v", test.num, test.name, res, test.out)
			}
		}
		ctrl.Finish()
	}

}

func TestCheck(t *testing.T) {

	usr := &user.User{Id: 1, Name: "alex", HashPass: goodHash}
	err := errors.New("wrong test")

	unitTests := []struct {
		num         int
		name        string
		getUser1    string
		getUser2    *user.User
		getUser3    error
		needGetUser bool
		check1      string
		check2      string
		check3      error
		out         *user.User
		err         error
		isErr       bool
	}{
		{
			num:         1,
			name:        "good way",
			getUser1:    "alex",
			getUser2:    usr,
			getUser3:    nil,
			needGetUser: true,
			check1:      "token",
			check2:      "alex",
			check3:      nil,
			out:         usr,
			err:         nil,
			isErr:       false,
		},
		{
			num:         2,
			name:        "wrong check",
			getUser1:    "alex",
			getUser2:    usr,
			getUser3:    nil,
			needGetUser: false,
			check1:      "token",
			check2:      "",
			check3:      err,
			out:         nil,
			err:         err,
			isErr:       true,
		},
		{
			num:         3,
			name:        "getUser error",
			getUser1:    "alex",
			getUser2:    nil,
			getUser3:    err,
			needGetUser: true,
			check1:      "token",
			check2:      "alex",
			check3:      nil,
			out:         nil,
			err:         err,
			isErr:       true,
		},
	}

	for _, test := range unitTests {
		ctrl := gomock.NewController(t)

		aRepo := ar.NewMockAuthRepoInterface(ctrl)
		aRepo.
			EXPECT().
			Check(test.check1).
			Return(test.check2, test.check3)

		uRepo := ur.NewMockUserRepoInterface(ctrl)
		if test.needGetUser {
			uRepo.
				EXPECT().
				GetUser(test.getUser1).
				Return(test.getUser2, test.getUser3)
		}

		au := NewAuthUsecase(uRepo, aRepo, salt)
		res, err := au.Check("token")

		if test.isErr {
			if err == nil {
				t.Errorf("%v %v - got:%v, want:%v", test.num, test.name, res, test.out)
			}
		} else {
			if fmt.Sprint(res) != fmt.Sprint(test.out) {
				t.Errorf("%v %v - got:%v, want:%v", test.num, test.name, res, test.out)
			}
		}
		ctrl.Finish()
	}

}
