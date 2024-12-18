package model

import (
	"gorm.io/gorm"
)

type ImageTable struct {
	gorm.Model
	ImageIdentifier string `json:"image_identifier" gorm:"index"`
	BaseName        string `json:"base_name" gorm:"index"`
	ImageType       string `json:"image_type"`
	ImageFileData   []byte `json:"image_file_data"`
}
