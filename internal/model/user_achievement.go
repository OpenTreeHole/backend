package model

import (
	"time"
)

// UserAchievement 用户成就关联表
type UserAchievement struct {
	// 用户 ID
	UserID int `json:"user_id" gorm:"primaryKey"`

	// 用户
	User *User `json:"user" gorm:"foreignKey:UserID"`

	// 成就 ID
	AchievementID int `json:"achievement_id" gorm:"primaryKey"`

	// 成就
	Achievement *Achievement `json:"achievement" gorm:"foreignKey:AchievementID"`

	// 获取日期
	ObtainDate time.Time `json:"obtain_date"`
}
