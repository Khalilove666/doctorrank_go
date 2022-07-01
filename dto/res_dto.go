package dto

import (
	"doctorrank_go/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginResDTO struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Email     string             `bson:"email" json:"email"`
	Username  string             `bson:"username" json:"username"`
	Token     string             `bson:"token" json:"token"`
	Role      string             `bson:"role" json:"role"`
	Img       string             `bson:"img" json:"img"`
	Contact   models.UserContact `bson:"contact" json:"contact"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}
