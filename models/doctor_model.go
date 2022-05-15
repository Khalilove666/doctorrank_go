package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Doctor struct {
	Id           primitive.ObjectID `bson:"_id" json:"_id"`
	UserId       primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title        string             `bson:"title" json:"title" validate:"required"`
	FirstName    string             `bson:"first_name" json:"first_name" validate:"required"`
	LastName     string             `bson:"last_name" json:"last_name" validate:"required"`
	Img          string             `bson:"img" json:"img"`
	About        string             `bson:"about" json:"about"`
	ProfessionId primitive.ObjectID `bson:"profession_id" json:"profession_id"`
	HospitalId   primitive.ObjectID `bson:"hospital_id" json:"hospital_id"`
	Experience   []Experience       `bson:"experience" json:"experience"`
	Education    []Education        `bson:"education" json:"education"`
	Contact      Contact            `bson:"contact" json:"contact"`
	CreatedAt    int64              `bson:"created_at" json:"created_at"`
	UpdatedAt    int64              `bson:"updated_at" json:"updated_at"`
}

type Experience struct {
	Profession string `bson:"profession" json:"profession"`
	Hospital   string `bson:"hospital" json:"hospital"`
	Field      string `bson:"field" json:"field"`
	TermStart  int64  `bson:"term_start" json:"term_start"`
	TermEnd    int64  `bson:"term_end" json:"term_end"`
	Country    string `bson:"country" json:"country"`
}

type Education struct {
	Degree      string `bson:"degree" json:"degree"`
	Major       string `bson:"major" json:"major"`
	Institution string `bson:"institution" json:"institution"`
	TermStart   int64  `bson:"term_start" json:"term_start"`
	TermEnd     int64  `bson:"term_end" json:"term_end"`
	Country     string `bson:"country" json:"country"`
}
type Contact struct {
	Phone    string `bson:"phone" json:"phone"`
	Email    string `bson:"email" json:"email"`
	Facebook string `bson:"facebook" json:"facebook"`
}
