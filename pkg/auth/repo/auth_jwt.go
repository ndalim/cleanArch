package repo

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"redditapp/pkg/user"
	"strconv"
	"time"
)

type jwtUser struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}

type jwtClaim struct {
	jwt.StandardClaims
	User *jwtUser `json:"user"`
}

type jwtResponse struct {
	Token string `json:"token"`
}

type AuthRepoMem struct {
	Secret string
	Ttl    int64
}

var (
	ErrorWrongSingIn = errors.New("bad login/password")
	ErrorSession     = errors.New("wrong session")
)

func NewAuthRepoMem(secret string, ttl int64) *AuthRepoMem {
	return &AuthRepoMem{
		Secret: secret,
		Ttl:    ttl,
	}
}

func (a *AuthRepoMem) Check(tokenString string) (string, error) {

	claims := &jwtClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(tok *jwt.Token) (interface{}, error) {
		return []byte(a.Secret), nil
	})

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", ErrorSession
	}

	return claims.User.Username, nil
}

func (a *AuthRepoMem) Session(usr *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwtClaim{
			User: &jwtUser{
				Username: usr.Name,
				Id:       strconv.Itoa(usr.Id),
			},
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + a.Ttl,
				IssuedAt:  time.Now().Unix(),
			},
		})

	sign, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", ErrorSession
	}

	return sign, nil
}
