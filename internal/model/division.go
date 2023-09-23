package model

import (
	"time"
)

type Division struct {
	// 分区 id
	ID int `json:"id" gorm:"primaryKey"`

	// 创建时间：只有管理员能创建分区
	CreatedAt time.Time `json:"time_created" gorm:"not null"`

	// 更新时间：只有管理员能更新分区，包括修改分区名称、分区详情、置顶的树洞
	UpdatedAt time.Time `json:"time_updated" gorm:"not null"`

	// 分区名称
	Name string `json:"name" gorm:"unique;size:10"`

	// 分区详情
	Description string `json:"description" gorm:"size:64"`

	// 置顶的树洞 id，按照顺序
	Pinned []int `json:"-" gorm:"serializer:json;not null;default:\"[]\""`
}
