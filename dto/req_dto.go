package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
)

type ImageDTO struct {
	File        *multipart.FileHeader `form:"file" validate:"file"`
	Coordinates struct {
		Left   int `bson:"left" json:"left" validate:"gte=0"`
		Top    int `bson:"top" json:"top" validate:"gte=0"`
		Width  int `bson:"width" json:"width" validate:"required,gt=0"`
		Height int `bson:"height" json:"height" validate:"required,gt=0"`
	} `form:"coordinates" validate:"json"`
}

type RegisterDTO struct {
	FirstName string `bson:"first_name" json:"first_name" validate:"required"`
	LastName  string `bson:"last_name" json:"last_name" validate:"required"`
	Email     string `bson:"email" json:"email" validate:"email,required"`
	Username  string `bson:"username" json:"username" validate:"required"`
	Password  string `bson:"password" json:"password" validate:"required,min=8"`
}

type LoginDTO struct {
	Login      string `bson:"login" json:"login" validate:"required"`
	Password   string `bson:"password" json:"password" validate:"required"`
	RememberMe bool   `bson:"remember_me" json:"remember_me" validate:"required"`
}

type PasswordDTO struct {
	OldPassword string `bson:"old_password" json:"old_password" validate:"required"`
	NewPassword string `bson:"new_password" json:"new_password" validate:"required,min=8"`
}

type PasswordResetDTO struct {
	NewPassword string `bson:"new_password" json:"new_password" validate:"required,min=8"`
}

type UserUpdateDTO struct {
	FieldName string `bson:"field_name" json:"field_name" validate:"required,oneof=first_name last_name contact_email contact_phone contact_facebook"`
	Value     string `bson:"value" json:"value" validate:"required"`
}

type DoctorUpdateDTO struct {
	FieldName string `bson:"field_name" json:"field_name" validate:"required,oneof=title first_name last_name about contact_email contact_phone contact_facebook profession_id hospital_id"`
	Value     string `bson:"value" json:"value" validate:"required"`
}

type DoctorExperienceUpdateDTO struct {
	Action string             `bson:"action" json:"action" validate:"required,oneof=create edit delete"`
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Value  struct {
		Profession string `bson:"profession" json:"profession"`
		Hospital   string `bson:"hospital" json:"hospital"`
		Field      string `bson:"field" json:"field"`
		TermStart  int64  `bson:"term_start" json:"term_start"`
		TermEnd    int64  `bson:"term_end" json:"term_end"`
		Country    string `bson:"country" json:"country"`
	} `bson:"value" json:"value" validate:"required"`
}

type DoctorEducationUpdateDTO struct {
	Action string             `bson:"action" json:"action" validate:"required,oneof=create edit delete"`
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	Value  struct {
		Degree      string `bson:"degree" json:"degree"`
		Major       string `bson:"major" json:"major"`
		Institution string `bson:"institution" json:"institution"`
		TermStart   int64  `bson:"term_start" json:"term_start"`
		TermEnd     int64  `bson:"term_end" json:"term_end"`
		Country     string `bson:"country" json:"country"`
	} `bson:"value" json:"value" validate:"required"`
}

type HospitalDTO struct {
	Name string `bson:"name" json:"name" validate:"required"`
}
