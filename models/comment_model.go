package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	DoctorId  string             `bson:"doctor_id" json:"doctor_id"`
	UserId    string             `bson:"user_id" json:"user_id"`
	Text      string             `bson:"text" json:"text"`
	Rate      float64            `bson:"rate" json:"rate"`
	Likes     []string           `bson:"likes" json:"likes"`
	Dislikes  []string           `bson:"dislikes" json:"dislikes"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}
