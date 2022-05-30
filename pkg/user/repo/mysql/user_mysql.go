package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"redditapp/pkg/user"
)

type UserRepoMysql struct {
	DB *sql.DB
}

func NewUserRepoMysql(db *sql.DB) UserRepoMysql {
	r := &UserRepoMysql{
		DB: db,
	}
	return *r
}

func (u UserRepoMysql) GetUser(name string) (*user.User, error) {

	rows, err := u.DB.Query("select id, name, hashpass from users where name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usr := &user.User{}
	var pass interface{}
	ok := false

	for rows.Next() {
		err = rows.Scan(&usr.Id, &usr.Name, &pass)

		if err != nil {
			return nil, err
		}
		if pass != nil {
			usr.HashPass = pass.([]byte)
		}
		ok = true
	}

	if ok {
		return usr, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func (u UserRepoMysql) GetAll() ([]*user.User, error) {

	res := make([]*user.User, 0)

	usr := &user.User{}
	var pass interface{}

	rows, err := u.DB.Query("select id, name, hashpass from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&usr.Id, &usr.Name, &pass)
		if err != nil {
			return nil, err
		}

		if pass != nil {
			usr.HashPass = pass.([]byte)
		}

		res = append(res, &user.User{
			Id:       usr.Id,
			Name:     usr.Name,
			HashPass: usr.HashPass,
		})
	}

	return res, nil
}

func (u UserRepoMysql) Create(name string, pass []byte) (*user.User, error) {

	usr, err := u.GetUser(name)
	if err == nil {
		fmt.Println("user_mysql Create: 1")
		return nil, user.ErrCreateUser
	}

	usr = &user.User{
		Name:     name,
		HashPass: pass,
	}

	tx, err := u.DB.Begin()
	if err != nil {
		fmt.Println("user_mysql Create: 2")
		return nil, user.ErrCreateUser
	}

	_, err = tx.Exec("insert into users (name, hashpass) values (?, ?)", name, pass)
	if err != nil {
		tx.Rollback()
		fmt.Println("user_mysql Create: 3 ", string(pass), err)
		return nil, user.ErrCreateUser
	}
	tx.Commit()

	usr, err = u.GetUser(name)
	if err != nil {
		tx.Rollback()
		fmt.Println("user_mysql Create: 4 ", err)
		return nil, user.ErrCreateUser
	}

	return usr, nil
}
