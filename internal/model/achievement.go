package model

import (
	"time"
)

// Achievement 成就
type Achievement struct {
	// 成就 ID , primary key
	ID int `json:"id"`

	// 创建时间
	CreatedAt time.Time `json:"created_at"`

	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`

	// 成就名称
	Name string `json:"name"`

	// 可能是成就作用域？
	Domain *string `json:"domain"`
}
