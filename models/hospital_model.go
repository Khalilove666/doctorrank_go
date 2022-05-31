package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hospital struct {
	Id   primitive.ObjectID `bson:"_id" json:"_id"`
	Name string             `bson:"name" json:"name"`
	Img  string             `bson:"img" json:"img"`
}
