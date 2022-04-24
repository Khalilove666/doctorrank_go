package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Doctor struct {
	Id         primitive.ObjectID `bson:"_id" json:"_id"`
	UserId     string             `bson:"user_id" json:"user_id"`
	Title      string             `bson:"title" json:"title" validate:"required"`
	FirstName  string             `bson:"first_name" json:"first_name" validate:"required"`
	LastName   string             `bson:"last_name" json:"last_name" validate:"required"`
	Img        string             `bson:"img" json:"img"`
	About      string             `bson:"about" json:"about"`
	Profession Profession         `bson:"profession" json:"profession"`
	Hospital   Hospital           `bson:"hospital" json:"hospital"`
	Experience []struct {
	} `bson:"experience" json:"experience"`
	Education []struct {
	} `bson:"education" json:"education"`
	Contact struct {
		Phone    string `bson:"phone" json:"phone"`
		Email    string `bson:"email" json:"email"`
		Facebook string `bson:"facebook" json:"facebook"`
	} `bson:"contact" json:"contact"`
	CreatedAt int64 `bson:"created_at" json:"created_at"`
	UpdatedAt int64 `bson:"updated_at" json:"updated_at"`
}

type CompactDoctor struct {
	Id         primitive.ObjectID `bson:"_id" json:"_id"`
	UserId     string             `bson:"user_id" json:"user_id"`
	Title      string             `bson:"title" json:"title" validate:"required"`
	FirstName  string             `bson:"first_name" json:"first_name" validate:"required"`
	LastName   string             `bson:"last_name" json:"last_name" validate:"required"`
	Img        string             `bson:"img" json:"img"`
	Profession Profession         `bson:"profession" json:"profession"`
	Hospital   Hospital           `bson:"hospital" json:"hospital"`
}
