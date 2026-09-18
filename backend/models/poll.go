package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Option struct {
	Text  string `bson:"text" json:"text"`
	Votes int64  `bson:"votes" json:"votes"`
}

type Poll struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Question  string        `bson:"question" json:"question"`
	Options   []Option      `bson:"options" json:"options"`
	CreatedBy bson.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
