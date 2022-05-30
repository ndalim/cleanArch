package mysql

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"redditapp/pkg/user"
	"testing"
)

var testPass = []byte("password")

func TestCreateUser(t *testing.T) {

	//errSql := errors.New("bad query")

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error("crush")
		return
	}

	repo := NewUserRepoMysql(db)

	queryText := "select id, name, hashpass from users where name = ?"

	s1 := mock.NewRows([]string{"id", "name", "hashpass"})
	mock.ExpectQuery(queryText).
		WithArgs("alex").
		WillReturnRows(s1)

	mock.ExpectBegin()
	r1 := sqlmock.NewResult(1, 1)
	mock.ExpectExec("insert into users (name, hashpass) values (?, ?)").
		WithArgs("alex", testPass).
		WillReturnResult(r1)

	s2 := mock.NewRows([]string{"id", "name", "hashpass"}).
		AddRow(1, "alex", "password")
	mock.ExpectQuery(queryText).
		WithArgs("alex").
		WillReturnRows(s2)
	mock.ExpectCommit()

	unitTests := []struct {
		num   int
		name  string
		in    string
		out   *user.User
		isErr bool
		err   error
	}{
		{
			num:  1,
			name: "good search",
			in:   "alex",
			out: &user.User{
				Id:       1,
				Name:     "alex",
				HashPass: testPass,
			},
			isErr: false,
			err:   nil,
		},
	}

	for _, test := range unitTests {
		res, err := repo.Create(test.in, testPass)
		//fmt.Println(res, err)
		if test.isErr {
			if err == nil {
				t.Errorf("#%v got:%v want:%v", test.num, false, true)
			}
		} else {
			if fmt.Sprint(res) != fmt.Sprint(test.out) {
				t.Errorf("#%v got:%v want:%v", test.num, res, test.out)
			}
		}
	}

}

func T_estCreateUser(t *testing.T) {

	//errSql := errors.New("bad query")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("crush")
		return
	}

	repo := NewUserRepoMysql(db)

	queryText := "select id, name, hashpass from users where name = ?"

	mock.ExpectBegin()

	s1 := mock.NewRows([]string{"id", "name", "hashpass"})
	mock.ExpectQuery(queryText).
		WithArgs("alex").
		WillReturnRows(s1)

	//mock.ExpectBegin()
	r1 := sqlmock.NewResult(1, 1)
	mock.ExpectExec("insert into users (name, hashpass) values (?, ?)").
		WithArgs("alex", testPass).
		WillReturnResult(r1)
	//mock.ExpectCommit()

	s2 := mock.NewRows([]string{"id", "name", "hashpass"}).
		AddRow(1, "john", "password")
	mock.ExpectQuery(queryText).
		WithArgs("john").
		WillReturnRows(s2)

	mock.ExpectBegin()
	r2 := sqlmock.NewResult(1, 1)
	mock.ExpectExec("insert into users (name, hashpass) values (?, ?)").
		WithArgs("john", testPass).
		WillReturnResult(r2)
	mock.ExpectRollback()

	unitTests := []struct {
		num   int
		name  string
		in    string
		out   *user.User
		isErr bool
		err   error
	}{
		{
			num:  1,
			name: "good search",
			in:   "alex",
			out: &user.User{
				Id:       1,
				Name:     "alex",
				HashPass: testPass,
			},
			isErr: false,
			err:   nil,
		},
		{
			num:   2,
			name:  "bad search",
			in:    "john",
			out:   nil,
			isErr: true,
			err:   user.ErrCreateUser,
		},
	}

	for _, test := range unitTests {
		res, err := repo.Create(test.in, []byte("password"))
		//fmt.Println(res, err)
		if test.isErr {
			if err == nil {
				t.Errorf("#%v got:%v want:%v", test.num, false, true)
			}
		} else {
			if fmt.Sprint(res) != fmt.Sprint(test.out) {
				t.Errorf("#%v got:%v want:%v", test.num, res, test.out)
			}
		}
	}

}

func TestGetUser(t *testing.T) {

	errSql := errors.New("bad query")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("crush")
		return
	}
	defer db.Close()

	queryText := "select id, name, hashpass from users where name = ?"
	repo := NewUserRepoMysql(db)

	r1 := mock.NewRows([]string{"id", "name", "hashpass"}).AddRow(1, "alex", "password")
	mock.ExpectQuery(queryText).WithArgs("alex").WillReturnRows(r1)

	r2 := mock.NewRows([]string{"id", "name", "hashpass"})
	mock.ExpectQuery(queryText).WithArgs("john").WillReturnRows(r2)

	r3 := mock.NewRows([]string{"id", "name"}).AddRow("3", "sam")
	mock.ExpectQuery(queryText).WithArgs("sam").WillReturnRows(r3)

	unitTests := []struct {
		num   int
		name  string
		in    string
		out   *user.User
		isErr bool
		err   error
	}{
		{
			num:  1,
			name: "good search",
			in:   "alex",
			out: &user.User{
				Id:       1,
				Name:     "alex",
				HashPass: testPass,
			},
			isErr: false,
			err:   nil,
		},
		{
			num:   2,
			name:  "bad search",
			in:    "john",
			out:   nil,
			isErr: true,
			err:   errSql,
		},
		{
			num:   3,
			name:  "error Query",
			in:    "sam",
			out:   nil,
			isErr: true,
			err:   errSql,
		},
		{
			num:   4,
			name:  "error Scan",
			in:    "val",
			out:   nil,
			isErr: true,
			err:   errSql,
		},
	}

	for _, test := range unitTests {
		res, err := repo.GetUser(test.in)
		if test.isErr {
			if err == nil {
				t.Errorf("#%v got:%v want:%v", test.num, false, true)
			}
		} else {
			if fmt.Sprint(res) != fmt.Sprint(test.out) {
				t.Errorf("#%v got:%v want:%v", test.num, res, test.out)
			}
		}
	}

}
