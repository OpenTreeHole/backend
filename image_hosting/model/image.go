package model

import (
	"gorm.io/gorm"
)

type ImageTable struct {
	gorm.Model
	ImageIdentifier string `json:"image_identifier" gorm:"uniqueIndex;size:20"`
	// ImageIdentifier  string `json:"image_identifier" gorm:"index:image_identifier_idx(20);unique"`
	OriginalFileName string `json:"original_file_name" gorm:"index"`
	ImageType        string `json:"image_type"`
	ImageFileData    []byte `json:"image_file_data"`
}
