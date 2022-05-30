package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"redditapp/pkg/auth/repo"
	mw "redditapp/pkg/middleware"
	ursql "redditapp/pkg/user/repo/mysql"

	"go.uber.org/zap"

	ud "redditapp/pkg/auth/delivery"
	cd "redditapp/pkg/comment/delivery"
	pd "redditapp/pkg/post/delivery"

	uuc "redditapp/pkg/auth/usecase"
	cuc "redditapp/pkg/comment/usecase"
	puc "redditapp/pkg/post/usecase"

	cr "redditapp/pkg/comment/repo"
	pr "redditapp/pkg/post/repo"
	//ur "redditapp/pkg/user/repo"

	_ "github.com/go-sql-driver/mysql"
)

var (
	AppStartError = errors.New("Application start error")
)

type configData struct {
	Auth_ttl        int64  `json:"auth_ttl"`
	Auth_salt       string `json:"auth_salt"`
	Auth_jwt_secret string `json:"auth_jwt_secret"`
	Redis_path      string `json:"redis_path"`
}

func getConfig() (*configData, error) {

	configPath := flag.String("config", "../../config/config.json", "server config")
	cf, err := os.Open(*configPath)

	if err != nil {
		return nil, err
	}

	cr := io.Reader(cf)
	confData, err := io.ReadAll(cr)
	if err != nil {
		return nil, err
	}

	conf := &configData{}
	err = json.Unmarshal(confData, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func NewLogger() (*zap.SugaredLogger, error) {

	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "../../logs/log"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, err
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil

}

func main() {

	conf, err := getConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger, err := NewLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()

	userDB, err := sql.Open("mysql", "root:123@tcp(localhost:3306)/reddit")
	defer userDB.Close()

	authDB, err := redis.DialURL(conf.Redis_path)
	defer authDB.Close()

	//userRepo := ur.NewUserRepoMem()
	userRepo := ursql.NewUserRepoMysql(userDB)

	postRepo := pr.NewPostRepoMem()
	commentRepo := cr.NewCommentRepoMem()

	authRepo := repo.NewAuthRepoRedis(authDB, conf.Auth_jwt_secret, conf.Auth_ttl)
	//authRepo := repo.NewAuthRepoMem(conf.Auth_jwt_secret, conf.Auth_ttl)

	auth := uuc.NewAuthUsecase(userRepo, authRepo, conf.Auth_salt)
	authHandle := ud.NewAuthHandler(auth, userRepo)

	comment := cuc.NewCommentUsecase(commentRepo)
	commentHandle := cd.NewCommentHandle(comment)

	post := puc.NewPostUsecase(postRepo)
	postHandle := pd.NewPostHandle(post, comment)

	//fmt.Println(authHandle, commentHandle, postHandle)
	//
	//basicMux := http.NewServeMux()
	//basicMux.Handle("/", http.FileServer(http.Dir("../../template")))
	//basicMux.Handle("/static/", http.FileServer(http.Dir("../../template")))

	basicMux := mux.NewRouter()
	basicMux.Handle("/", http.FileServer(http.Dir("../../template")))
	basicMux.Handle("/static", http.FileServer(http.Dir("../../template")))

	basicMux.HandleFunc("/api/login", authHandle.SignIn).Methods("POST")
	basicMux.HandleFunc("/api/register", authHandle.SignUp).Methods("POST")

	basicMux.HandleFunc("/api/user/{login}", postHandle.GetByUser).Methods("GET")

	basicMux.HandleFunc("/api/posts", mw.Auth(*auth, postHandle.AddPost)).Methods("POST")
	basicMux.HandleFunc("/api/posts/", postHandle.GetAll).Methods("GET")
	basicMux.HandleFunc("/api/posts/{category_name}", postHandle.GetByCategory).Methods("GET")

	basicMux.HandleFunc("/api/post/{id}", postHandle.GetPost).Methods("GET")
	basicMux.HandleFunc("/api/post/{id}", mw.WrapHeaders(mw.Auth(*auth, postHandle.AddComment))).Methods("POST")
	basicMux.HandleFunc("/api/post/{id}", mw.Auth(*auth, postHandle.Delete)).Methods("DELETE")

	basicMux.HandleFunc("/api/post/{id}/upvote", mw.Auth(*auth, postHandle.UpVote)).Methods("GET")
	basicMux.HandleFunc("/api/post/{id}/downvote", mw.Auth(*auth, postHandle.DownVote)).Methods("GET")
	basicMux.HandleFunc("/api/post/{id}/unvote", mw.Auth(*auth, postHandle.UnVote)).Methods("GET")
	basicMux.HandleFunc("/api/post/{id}/{comment_id}", mw.Auth(*auth, commentHandle.Delete)).Methods("DELETE")

	basicMux.HandleFunc("/api/test/user/{name}", authHandle.GetUser).Methods("GET")
	basicMux.HandleFunc("/api/test/users", authHandle.GetAll).Methods("GET")

	siteMux := mw.AccessLog(logger, basicMux)
	siteMux = mw.PanicMiddleware(siteMux)

	logger.Infow("server :8086 has started")
	logger.Fatal(http.ListenAndServe(":8086", siteMux))

}
