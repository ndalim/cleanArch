package repo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	"redditapp/pkg/post"
	"redditapp/pkg/user"
	"redditapp/tools"
	"time"
)

var (
	ErrGetPost = errors.New("wrong getting post")
	ErrAddPost = errors.New("wrong adding post")
)

var Now = time.Now

var CreateID = primitive.NewObjectID

func postMongo2Model(pst bson.M) *post.Post {

	aBson := pst["author"].(bson.M)
	a := &user.User{
		Id:   int(aBson["id"].(int32)),
		Name: aBson["name"].(string),
	}

	vBson := pst["votes"].(bson.A)
	votes := []post.Vote{}
	for _, v := range vBson {
		voteBson := v.(bson.M)
		vote := post.Vote{
			User: voteBson["user"].(string),
			Vote: int(voteBson["vote"].(int32)),
		}
		votes = append(votes, vote)

	}

	res := &post.Post{
		Score:            int(pst["score"].(int32)),
		Views:            int(pst["views"].(int32)),
		TypePost:         post.PType(pst["typePost"].(string)),
		Title:            pst["title"].(string),
		Url:              pst["url"].(string),
		Author:           a,
		Category:         post.PCategory(pst["category"].(string)),
		Text:             pst["text"].(string),
		Votes:            votes,
		Created:          pst["created"].(string),
		UpvotePercentage: 0,
		Id:               pst["_id"].(primitive.ObjectID).Hex(),
	}
	tools.CalculateUpvotePercentage(res)
	return res
}

func postModel2Mongo(pst *post.Post) bson.M {

	if pst.Id == "" {
		pst.Id = CreateID().Hex()
	}
	objId, _ := primitive.ObjectIDFromHex(pst.Id)

	a := bson.M{
		"id":   pst.Author.Id,
		"name": pst.Author.Name,
	}

	votes := make(bson.A, len(pst.Votes))
	for i, v := range pst.Votes {
		vote := bson.M{
			"user": v.User,
			"vote": v.Vote,
		}
		votes[i] = vote
	}

	res := bson.M{
		"score":            pst.Score,
		"views":            pst.Views,
		"typePost":         pst.TypePost,
		"title":            pst.Title,
		"url":              pst.Url,
		"author":           a,
		"category":         pst.Category,
		"text":             pst.Text,
		"votes":            votes,
		"created":          pst.Created,
		"upvotePercentage": pst.UpvotePercentage,
		"_id":              objId,
	}
	return res
}

type PostRepoMongo struct {
	db    *mongo.Database
	posts *mongo.Collection
	votes *mongo.Collection
}

func NewPostRepoMongo(db *mongo.Database) *PostRepoMongo {
	p := &PostRepoMongo{
		db:    db,
		posts: db.Collection("posts"),
		votes: db.Collection("votes"),
	}

	return p
}

func (r *PostRepoMongo) GetAll() ([]post.Post, error) {

	res := make([]post.Post, 0)

	ctx := context.TODO()

	cur, err := r.posts.Find(ctx, bson.M{})
	if err != nil {
		return nil, ErrGetPost
	}

	for cur.Next(ctx) {
		pst := bson.M{}
		err := cur.Decode(pst)
		if err != nil {
			return nil, ErrGetPost
		}
		res = append(res, *postMongo2Model(pst))
	}

	return res, nil
}

func (r *PostRepoMongo) GetByCategory(c post.PCategory) ([]post.Post, error) {

	res := make([]post.Post, 0)

	ctx := context.TODO()
	cur, err := r.posts.Find(ctx, bson.M{"category": c})
	if err != nil {
		return nil, ErrGetPost
	}

	for cur.Next(ctx) {
		pst := bson.M{}
		err := cur.Decode(pst)
		if err != nil {
			return nil, ErrGetPost
		}
		res = append(res, *postMongo2Model(pst))
	}

	return res, nil
}

func (r *PostRepoMongo) GetByUser(user string) ([]post.Post, error) {

	res := make([]post.Post, 0)

	ctx := context.TODO()
	cur, err := r.posts.Find(ctx, bson.M{"author.name": user})
	if err != nil {
		return nil, ErrGetPost
	}

	for cur.Next(ctx) {
		pst := bson.M{}
		err := cur.Decode(pst)
		if err != nil {
			return nil, ErrGetPost
		}
		res = append(res, *postMongo2Model(pst))
	}

	return res, nil
}

func (r *PostRepoMongo) GetPost(id string) (*post.Post, error) {

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrGetPost
	}

	ctx := context.TODO()

	pst := bson.M{}
	res := r.posts.FindOne(ctx, bson.M{"_id": objId})
	err = res.Decode(pst)
	if err != nil {
		return nil, ErrGetPost
	}
	return postMongo2Model(pst), nil

}

func (r *PostRepoMongo) ShowPost(id string) (*post.Post, error) {

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrGetPost
	}

	ctx := context.TODO()

	pst := bson.M{}
	res := r.posts.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$inc": bson.M{"views": 1}})
	err = res.Decode(pst)
	if err != nil {
		return nil, ErrGetPost
	}
	return postMongo2Model(pst), nil
}

func (r *PostRepoMongo) AddPost(pst *post.Post, author *user.User) (string, error) {

	pst.Created = Now().Format(time.RFC3339)
	pst.Votes = make([]post.Vote, 0)
	pst.Author = author

	ctx := context.TODO()
	pstBson := postModel2Mongo(pst)
	fmt.Println(pstBson)
	_, err := r.posts.InsertOne(ctx, pstBson)
	if err != nil {
		//return "", ErrAddPost
		return "", err
	}

	return pst.Id, nil
}

func (r *PostRepoMongo) SetVote(id string, user string, vote int) error {

	objId, err := primitive.ObjectIDFromHex(id)
	ctx := context.TODO()

	qFind := bson.M{"_id": objId}
	qDelete := bson.M{
		"$pull": bson.M{
			"votes": bson.M{
				"user": user,
			},
		},
	}
	qInsert := bson.M{
		"$push": bson.M{
			"votes": bson.M{
				"user": user,
				"vote": vote,
			},
		},
	}

	commands := []mongo.WriteModel{
		mongo.NewUpdateOneModel().SetFilter(qFind).SetUpdate(qDelete),
	}

	if vote != 0 {
		commands = append(commands, mongo.NewUpdateOneModel().SetFilter(qFind).SetUpdate(qInsert))
	}

	_, err = r.posts.BulkWrite(ctx, commands)
	if err != nil {
		return err
	}

	return nil

}

func (r *PostRepoMongo) UpVote(id string, user string) error {
	return r.SetVote(id, user, 1)
}

func (r *PostRepoMongo) DownVote(id string, user string) error {
	return r.SetVote(id, user, -1)
}

func (r *PostRepoMongo) UnVote(id string, user string) error {
	return r.SetVote(id, user, 0)
}

func (r *PostRepoMongo) Delete(id string, user string) error {

	objId, _ := primitive.ObjectIDFromHex(id)
	ctx := context.TODO()

	_, err := r.posts.DeleteOne(ctx, bson.M{"_id": objId, "author.name": user})
	if err != nil {
		return err
	}
	return nil
}
