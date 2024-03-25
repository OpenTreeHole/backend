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

// UserAchievement 用户成就关联表
type UserAchievement struct {
	UserID        int       `json:"user_id" gorm:"primaryKey"`        // 用户 ID
	AchievementID int       `json:"achievement_id" gorm:"primaryKey"` // 成就 ID
	ObtainDate    time.Time `json:"obtain_date"`                      // 获得日期
}
