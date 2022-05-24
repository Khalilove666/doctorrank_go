package dto

import "mime/multipart"

type ImageDTO struct {
	File        *multipart.FileHeader `form:"file" validate:"file"`
	Coordinates struct {
		Left   int `bson:"left" json:"left" validate:"gte=0"`
		Top    int `bson:"top" json:"top" validate:"gte=0"`
		Width  int `bson:"width" json:"width" validate:"required,gt=0"`
		Height int `bson:"height" json:"height" validate:"required,gt=0"`
	} `form:"coordinates" validate:"json"`
}
