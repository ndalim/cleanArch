package post

import (
	"errors"
	"redditapp/pkg/user"
)

type PType string

type Vote struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

type PCategory string

var ErrorPostNotFound = errors.New("post not found")

const (
	typeLink PType = "link"
	typeText PType = "text"
)

type Post struct {
	Score            int        `json:"score"`
	Views            int        `json:"views"`
	TypePost         PType      `json:"type"`
	Title            string     `json:"title"`
	Url              string     `json:"url"`
	Author           *user.User `json:"author"`
	Category         PCategory  `json:"category"`
	Text             string     `json:"text"`
	Votes            []Vote     `json:"votes"`
	Created          string     `json:"created"`
	UpvotePercentage int        `json:"upvotePercentage"`
	Id               string     `json:"id"`
}
