package model

import (
	"golang.org/x/sys/windows"
	"gorm.io/gorm"
	"time"
)

type ImageTable struct {
	gorm.Model
	BaseName      string       `json:"base_name"`
	SavingTime    time.Time    `json:"saving_time"`
	ImageType     string       `json:"image_type"`
	ImageFileData windows.BLOB `json:"image_file_data"`
}
