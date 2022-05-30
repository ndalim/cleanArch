package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"redditapp/pkg/auth/repo"
	mw "redditapp/pkg/middleware"
	"reflect"
	"strings"
	"testing"
	"time"

	ursql "redditapp/pkg/user/repo/mysql"

	ud "redditapp/pkg/auth/delivery"
	cd "redditapp/pkg/comment/delivery"
	pd "redditapp/pkg/post/delivery"

	uuc "redditapp/pkg/auth/usecase"
	cuc "redditapp/pkg/comment/usecase"
	puc "redditapp/pkg/post/usecase"

	cr "redditapp/pkg/comment/repo"
	pr "redditapp/pkg/post/repo"
)

var (
	client = http.Client{Timeout: 1 * time.Second}
)

func PrepareTest(db *sql.DB, dbMongo *mongo.Database) {
	command_list := []string{
		`DROP TABLE IF EXISTS users;`,
		`CREATE TABLE users (
  			id int NOT NULL AUTO_INCREMENT,
  			name varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
			hashpass blob,
  			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	for _, v := range command_list {
		_, err := db.Exec(v)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.TODO()
	posts := dbMongo.Collection("posts")
	posts.DeleteMany(ctx, bson.M{})

}

func CleanTest(db *sql.DB) {
	command_list := []string{
		`DROP TABLE IF EXISTS users;`,
	}

	for _, v := range command_list {
		_, err := db.Exec(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type UnitTest struct {
	num       int
	name      string
	method    string
	query     string
	body      string
	headAuth  string
	resStatus int
	resBody   string
}

func TestHandlers(t *testing.T) {

	db, err := sql.Open("mysql", "root:123@tcp(localhost:3306)/reddit")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := NewLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()

	authDB, err := redis.DialURL("redis://localhost:6379")
	defer authDB.Close()

	userRepo := ursql.NewUserRepoMysql(db)

	opt := options.Client().ApplyURI("mongodb://localhost:27017")
	clientMongo, err := mongo.NewClient(opt)
	if err != nil {
		fmt.Println("crush mongo create client")
		return
	}

	ctx := context.TODO()
	err = clientMongo.Connect(ctx)
	if err != nil {
		fmt.Println("crush mongo connect")
		return
	}

	err = clientMongo.Ping(ctx, nil)
	if err != nil {
		fmt.Println("crush mongo ping")
		return
	}

	dbMongo := clientMongo.Database("reddit")

	PrepareTest(db, dbMongo)
	//defer CleanTest(db)

	postRepo := pr.NewPostRepoMongo(dbMongo)
	commentRepo := cr.NewCommentRepoMongo(dbMongo)

	now := func() time.Time {
		return time.Date(2022, 5, 21, 0, 0, 0, 0, time.UTC)
	}

	repo.Now = now
	pr.Now = now

	repo.CreateID = func() string {
		return "WhSfyuRl"
	}
	pr.CreateID = func() primitive.ObjectID {
		res, _ := primitive.ObjectIDFromHex("626674010f5630cbfaed53f6")
		return res
	}

	authRepo := repo.NewAuthRepoRedis(authDB, "reddit_jwt_secret", 31*24*3600)

	auth := uuc.NewAuthUsecase(userRepo, authRepo, "reddit_pass_salt")
	authHandle := ud.NewAuthHandler(auth, userRepo)

	comment := cuc.NewCommentUsecase(commentRepo)
	commentHandle := cd.NewCommentHandle(comment)

	post := puc.NewPostUsecase(postRepo)
	postHandle := pd.NewPostHandle(post, comment)

	basicMux := mux.NewRouter()
	//basicMux.Handle("/", http.FileServer(http.Dir("../../template")))
	//basicMux.Handle("/static", http.FileServer(http.Dir("../../template")))

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

	ts := httptest.NewServer(siteMux)

	goodToken := `eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTU3Njk2MDAsImlhdCI6MTY1MzA5MTIwMCwidXNlciI6eyJ1c2VybmFtZSI6ImFzZGYxMjM0IiwiaWQiOiIxIn0sIklkIjoiV2hTZnl1UmwifQ.hplVRgneUJxIYmHJy11k79jO4OrcoMpLWpdp4Rf8sEdtghps5AIfl9CeHpVhvK_rZPxsB1L2Wd-oFbhy_DgqQQ`
	cases := []UnitTest{
		{
			num:       1,
			name:      "register",
			method:    "POST",
			query:     "/api/register",
			body:      `{"username":"asdf1234","password":"asdfasdf"}`,
			resStatus: 200,
			resBody:   fmt.Sprintf(`{"token":"%v"}`, goodToken),
		},
		{
			num:       2,
			name:      "login",
			method:    "POST",
			query:     "/api/login",
			body:      `{"username":"asdf1234","password":"asdfasdf"}`,
			resStatus: 200,
			resBody:   fmt.Sprintf(`{"token":"%v"}`, goodToken),
		},
		{
			num:       3,
			name:      "add post",
			method:    "POST",
			query:     "/api/posts",
			body:      `{"category":"music","type":"text","title":"123","text":"text 123"}`,
			headAuth:  "Bearer " + goodToken,
			resStatus: 201,
			resBody: `{
							"score": 0,
							"views": 0,
							"type": "text",
							"title": "123",
							"author": {
								"username": "asdf1234",
								"id": 1
							},
							"category": "music",
							"text": "text 123",
							"votes": [],
							"url": "",
							"created": "2022-05-21T00:00:00Z",
							"upvotePercentage": 0,
							"id": "626674010f5630cbfaed53f6"
						}`,
		},
		{
			num:       4,
			name:      "downvote",
			method:    "GET",
			query:     "/api/post/626674010f5630cbfaed53f6/downvote",
			body:      ``,
			headAuth:  "Bearer " + goodToken,
			resStatus: 200,
			resBody: `{
							"score": -1,
							"views": 1,
							"type": "text",
							"title": "123x2",
							"author": {
								"username": "asdf1234",
								"id": 1
							},
							"category": "music",
							"text": "123 123",
							"votes": [
								{
									"user": 1,
									"vote": -1
								}
							],
							"comments": [],
							"created": "2022-05-23T18:42:26.195Z",
							"upvotePercentage": 0,
							"id": "626674010f5630cbfaed53f6"
						}`,
		},
	}

	runTests(t, ts, cases)
}

func runTests(t *testing.T, ts *httptest.Server, cases []UnitTest) {
	for _, unit := range cases {

		body := strings.NewReader(unit.body)
		r, err := http.NewRequest(unit.method, ts.URL+unit.query, body)

		if unit.headAuth != "" {
			r.Header["Authorization"] = []string{unit.headAuth}
		}

		res, err := client.Do(r)

		resBody, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		if err != nil {
			t.Fatalf("%v", err)
			return
		}

		if res.StatusCode != unit.resStatus {
			t.Errorf("%v %v - wrong status, want: %v, got: %v", unit.num, unit.name, unit.resStatus, res.StatusCode)
			continue
		}

		want := []byte(unit.resBody)
		var objWant, objGot interface{}
		json.Unmarshal(resBody, &objGot)
		json.Unmarshal(want, &objWant)

		if !reflect.DeepEqual(objGot, objWant) {
			//fmt.Println(objGot, objWant)
			fmt.Println(string(resBody))
			t.Errorf("%v %v - wrond body, want: %v, got: %v", unit.num, unit.name, string(want), string(resBody))
			continue
		}
	}
}
