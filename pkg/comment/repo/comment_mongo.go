package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"redditapp/pkg/comment"
	"redditapp/pkg/user"
	"time"
)

type CommentRepoMongo struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func commentMongo2Model(comm bson.M) *comment.Comment {

	aBson := comm["author"].(bson.M)
	usr := &user.User{
		Id:   int(aBson["id"].(int32)),
		Name: aBson["name"].(string),
	}

	res := &comment.Comment{
		Created: comm["created"].(string),
		Author:  usr,
		Body:    comm["body"].(string),
		Id:      comm["_id"].(primitive.ObjectID).Hex(),
		PostId:  comm["postId"].(primitive.ObjectID).Hex(),
	}

	return res
}

func commentModel2Mongo(comm *comment.Comment) bson.M {

	if comm.Id == "" {
		comm.Id = primitive.NewObjectID().Hex()
	}
	objId, _ := primitive.ObjectIDFromHex(comm.Id)

	postId, _ := primitive.ObjectIDFromHex(comm.PostId)

	a := bson.M{
		"id":   comm.Author.Id,
		"name": comm.Author.Name,
	}

	res := bson.M{
		"created": comm.Created,
		"author":  a,
		"body":    comm.Body,
		"_id":     objId,
		"postId":  postId,
	}
	return res
}

func NewCommentRepoMongo(db *mongo.Database) *CommentRepoMongo {
	return &CommentRepoMongo{
		db:   db,
		coll: db.Collection("comments"),
	}
}

func (r *CommentRepoMongo) GetByPost(id string) ([]*comment.Comment, error) {

	postId, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	cur, err := r.coll.Find(ctx, bson.M{"postId": postId})
	if err != nil {
		return nil, err
	}

	res := make([]*comment.Comment, 0, 5)

	for cur.Next(ctx) {
		comm := bson.M{}
		err = cur.Decode(comm)
		if err != nil {
			return nil, err
		}
		res = append(res, commentMongo2Model(comm))
	}

	return res, nil
}

func (r *CommentRepoMongo) AddComment(comm comment.Comment) (string, error) {

	ctx := context.TODO()

	cBson := commentModel2Mongo(&comm)
	_, err := r.coll.InsertOne(ctx, cBson)
	if err != nil {
		return "", err
	}
	return comm.Id, nil
}

func (r *CommentRepoMongo) Delete(pstId, id, user string) error {

	objId, _ := primitive.ObjectIDFromHex(id)
	postId, _ := primitive.ObjectIDFromHex(pstId)

	ctx := context.TODO()
	_, err := r.coll.DeleteOne(ctx, bson.M{
		"_id":         objId,
		"postId":      postId,
		"author.name": user,
	})

	if err != nil {
		return err
	}
	return nil

}
