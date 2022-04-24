package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Profession struct {
	Id   primitive.ObjectID `bson:"_id" json:"_id"`
	Name string             `bson:"name" json:"name" validate:"required"`
}
