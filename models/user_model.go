package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName      string             `bson:"first_name" json:"first_name" validate:"required"`
	LastName       string             `bson:"last_name" json:"last_name" validate:"required"`
	Email          string             `bson:"email" json:"email" validate:"email,required"`
	Username       string             `bson:"username" json:"username" validate:"required"`
	Password       string             `bson:"password" json:"password" validate:"required,min=6"`
	Role           string             `bson:"role" json:"role"`
	Img            string             `bson:"img" json:"img"`
	EmailConfirmed bool               `bson:"email_confirmed" json:"email_confirmed"`
	Contact        UserContact        `bson:"contact" json:"contact"`
	CreatedAt      int64              `bson:"created_at" json:"created_at"`
	UpdatedAt      int64              `bson:"updated_at" json:"updated_at"`
}

type UserContact struct {
	Phone    string `bson:"phone" json:"phone"`
	Email    string `bson:"email" json:"email"`
	Facebook string `bson:"facebook" json:"facebook"`
}
