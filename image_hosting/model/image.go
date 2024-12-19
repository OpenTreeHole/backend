package model

import (
	"gorm.io/gorm"
)

type ImageTable struct {
	gorm.Model
	ImageIdentifier  string `json:"image_identifier" gorm:"uniqueIndex;size:20"`
	OriginalFileName string `json:"original_file_name" gorm:"index"`
	ImageType        string `json:"image_type"`
	ImageFileData    []byte `json:"image_file_data"`
}
