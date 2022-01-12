package model

import (
	"errors"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrNoRecord = errors.New("models: no matching record found") 

type Article struct {
	ID  primitive.ObjectID `json:"_id" bson:"_id"` 
	Title string `json:"title" bson:"title,omitempty"`
	Content string `json:"content" bson:"content,omitempty"`
	TimeStamp time.Time `json:"time_stamp" bson:"time_stamp,omitempty"`
}

func NewArticle() *Article {
	_article := new(Article)
	_article.ID = primitive.NewObjectID()
	_article.Title = ""
	_article.Content = ""
	_article.TimeStamp = time.Now()

	return _article
}