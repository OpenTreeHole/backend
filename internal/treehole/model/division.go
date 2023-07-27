package model

import (
	"gorm.io/datatypes"
)

type Division struct {
	ID          int
	Name        string
	Description string
	Pinned      datatypes.JSONSlice[int]
}
