package model

import (
	"time"
)

// ReviewRank 评教分数
type ReviewRank struct {
	Overall    int `json:"overall"`
	Content    int `json:"content"`    // 内容、风格方面
	Workload   int `json:"workload"`   // 工作量方面
	Assessment int `json:"assessment"` // 考核方面
}

// Review 评教
type Review struct {
	ID               int                `json:"id"`
	CreatedAt        time.Time          `json:"created_at" gorm:"not null"`
	UpdatedAt        time.Time          `json:"updated_at" gorm:"not null"`
	CourseID         int                `json:"course_id" gorm:"not null;index"`
	Course           *Course            `json:"course"`
	Title            string             `json:"title" gorm:"not null"`
	Content          string             `json:"content" gorm:"not null"`
	ReviewerID       int                `json:"reviewer_id" gorm:"not null;index"`
	Rank             *ReviewRank        `json:"rank" gorm:"embedded;embeddedPrefix:rank_"`
	UpvoteCount      int                `json:"upvote_count"`
	DownvoteCount    int                `json:"downvote_count"`
	ModifyCount      int                `json:"modify_count" gorm:"not null;default:0"`
	History          []*ReviewHistory   `json:"history"`
	Vote             []*ReviewVote      `json:"vote"`
	UserAchievements []*UserAchievement `json:"achievements" gorm:"foreignKey:UserID;references:ReviewerID"`
}

// ReviewHistory 评教修改历史
type ReviewHistory struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"` // 创建时间，原本是 time_created
	UpdatedAt time.Time `json:"updated_at"` // 更新时间，原本是 time_updated
	ReviewID  int       `json:"review_id" gorm:"not null;index"`
	AlterBy   int       `json:"alter_by" gorm:"not null"` // 修改人 ID
	Title     string    `json:"title" gorm:"not null"`    // 修改前的标题
	Content   string    `json:"content" gorm:"not null"`  // 修改前的内容
}

// ReviewVote 评教点赞/点踩详情
type ReviewVote struct {
	UserID   int `json:"user_id" gorm:"primaryKey"`   // 点赞或点踩人的 ID
	ReviewID int `json:"review_id" gorm:"primaryKey"` // 评教 ID
	Data     int `json:"data"`                        // 1 为点赞，-1 为点踩
}
