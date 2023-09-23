package model

import (
	"time"
)

// UserAchievement 用户成就关联表
type UserAchievement struct {
	// 用户 ID
	UserID int `json:"user_id"`

	// 成就 ID
	AchievementID int `json:"achievement_id"`

	// 获取日期
	ObtainDate time.Time `json:"obtain_date"`
}
