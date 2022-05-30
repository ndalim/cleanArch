package repo

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"redditapp/pkg/user"
	"strconv"
	"time"
)

var alphabet = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")

var Now = time.Now

var CreateID = func() string {
	res := make([]rune, 8)
	for i := range res {
		res[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(res)
}

type UserRedis struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}

type ClaimRedis struct {
	jwt.StandardClaims
	User *jwtUser `json:"user"`
	Id   string
}

type ResponseRedis struct {
	Token string `json:"token"`
}

type AuthRepoRedis struct {
	Conn   redis.Conn
	Secret string
	Ttl    int64
}

var (
	ErrorWrongSingInRedis = errors.New("bad login/password")
	ErrorSessionRedis     = errors.New("wrong session")
)

func NewAuthRepoRedis(conn redis.Conn, secret string, ttl int64) *AuthRepoRedis {
	return &AuthRepoRedis{
		Conn:   conn,
		Secret: secret,
		Ttl:    ttl,
	}
}

func (a *AuthRepoRedis) Check(tokenString string) (string, error) {

	claims := &ClaimRedis{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(tok *jwt.Token) (interface{}, error) {
		return []byte(a.Secret), nil
	})

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", ErrorSessionRedis
	}

	_, err = redis.String(a.Conn.Do("GET", claims.Id))
	if err != nil {
		return "", ErrorSessionRedis
	}

	return claims.User.Username, nil
}

func (a *AuthRepoRedis) Session(usr *user.User) (string, error) {
	id := CreateID()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		ClaimRedis{
			User: &jwtUser{
				Username: usr.Name,
				Id:       strconv.Itoa(usr.Id),
			},
			Id: id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: Now().Unix() + a.Ttl,
				IssuedAt:  Now().Unix(),
			},
		})

	sign, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", ErrorSessionRedis
	}

	_, err = a.Conn.Do("SET", id, sign)
	if err != nil {
		return "", ErrorSessionRedis
	}

	fmt.Println("auth_redis SESSION: ", id)
	return sign, nil
}


