package model

import (
	"time"
)

type ReviewRank struct {
	Overall int `json:"overall"`

	// 内容、风格方面
	Content int `json:"content"`

	// 工作量方面
	Workload int `json:"workload"`

	// 考核方面
	Assessment int `json:"assessment"`
}

// Review 评教
type Review struct {
	// 评教 ID , primary key
	ID int `json:"id"`

	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"not null"`

	// 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`

	// 课程 ID
	CourseID int     `json:"course_id" gorm:"not null;index"`
	Course   *Course `json:"course"`

	// 标题
	Title string `json:"title" gorm:"not null"`

	// 内容
	Content string `json:"content" gorm:"not null"`

	// 评教者
	ReviewerID int `json:"reviewer_id" gorm:"not null;index"`

	// 评分
	Rank *ReviewRank `json:"rank" gorm:"embedded;embeddedPrefix:rank_"`

	// 点赞数
	UpvoteCount int `json:"upvote_count"`

	// 点踩数
	DownvoteCount int `json:"downvote_count"`

	// 评教修改历史
	History []*ReviewHistory `json:"history"`

	// 评教点赞/点踩详情
	Vote []*ReviewVote `json:"vote"`

	// 用户成就
	UserAchievements []*UserAchievement `json:"achievements" gorm:"foreignKey:UserID;references:ReviewerID"`
}

// ReviewHistory 评教修改历史
type ReviewHistory struct {
	// 评教修改历史 ID , primary key
	ID int `json:"id"`

	// 创建时间，原本是 time_created
	CreatedAt time.Time `json:"created_at"`

	// 更新时间，原本是 time_updated
	UpdatedAt time.Time `json:"updated_at"`

	// 评教 ID
	ReviewID int `json:"review_id" gorm:"not null;index"`

	// 修改人 ID
	AlterBy int `json:"alter_by" gorm:"not null"`

	// 修改前的标题
	Title string `json:"title" gorm:"not null"`

	// 修改前的内容
	Content string `json:"content" gorm:"not null"`
}

// ReviewVote 评教点赞/点踩详情
type ReviewVote struct {
	// 点赞或点踩人的 ID
	UserID int `json:"user_id" gorm:"primaryKey"`

	// 评教 ID
	ReviewID int `json:"review_id" gorm:"primaryKey"`

	// 1 为点赞，-1 为点踩
	Data int `json:"data"`
}
