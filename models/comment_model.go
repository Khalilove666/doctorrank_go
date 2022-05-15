package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	DoctorId  primitive.ObjectID `bson:"doctor_id" json:"doctor_id"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Text      string             `bson:"text" json:"text" validate:"required"`
	Rate      float64            `bson:"rate" json:"rate" validate:"required,min=1,max=5"`
	Likes     []Like             `bson:"likes" json:"likes"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}

type Like struct {
	UserId primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status bool               `bson:"status" json:"status"`
}
