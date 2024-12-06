package model

import (
	"gorm.io/gorm"
)

type ImageTable struct {
	gorm.Model
	BaseName      string `json:"base_name"`
	ImageType     string `json:"image_type"`
	ImageFileData []byte `json:"image_file_data"`
}
