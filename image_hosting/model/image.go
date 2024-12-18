package model

import (
	"gorm.io/gorm"
)

type ImageTable struct {
	gorm.Model
	ImageIdentifier string `json:"image_identifier"`
	BaseName        string `json:"base_name" gorm:"index"`
	ImageType       string `json:"image_type" gorm:"index:idx_image_type"`
	ImageFileData   []byte `json:"image_file_data"`
}
