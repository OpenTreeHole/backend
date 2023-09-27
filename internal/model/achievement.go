package model

import (
	"time"
)

// Achievement 成就
type Achievement struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" gorm:"not null"` // 成就名称
	Domain    string    `json:"domain"`               // 可能是成就作用域？
}
