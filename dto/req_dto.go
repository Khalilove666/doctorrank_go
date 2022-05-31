package dto

import (
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
	Password  string `bson:"password" json:"password" validate:"required,min=6"`
}

type LoginDTO struct {
	Login      string `bson:"login" json:"login" validate:"required"`
	Password   string `bson:"password" json:"password" validate:"required"`
	RememberMe bool   `bson:"remember_me" json:"remember_me" validate:"required"`
}

type PasswordDTO struct {
	OldPassword string `bson:"old_password" json:"old_password" validate:"required"`
	NewPassword string `bson:"new_password" json:"new_password" validate:"required,min=6"`
}

type UserUpdateDTO struct {
	FirstName *string `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  *string `bson:"last_name,omitempty" json:"last_name,omitempty"`
}
