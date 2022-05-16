package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName string             `bson:"first_name" json:"first_name" validate:"required"`
	LastName  string             `bson:"last_name" json:"last_name" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"email,required"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	Password  string             `bson:"password" json:"password" validate:"required"`
	Role      string             `bson:"role" json:"role"`
	Img       string             `bson:"img" json:"img"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}

type LoginCredentials struct {
	Login      string `bson:"login" json:"login" validate:"required"`
	Password   string `bson:"password" json:"password" validate:"required"`
	RememberMe bool   `bson:"remember_me" json:"remember_me"`
}
