package model

import (
	"time"
)

// Review 评教
type Review struct {
	// 评教 ID , primary key
	ID int `json:"id"`

	// 创建时间
	CreatedAt time.Time `json:"created_at"`

	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`

	// 课程 ID
	CourseID int `json:"course_id"`

	// 标题
	Title string `json:"title"`

	// 内容
	Content string `json:"content"`

	// 评教时间
	ReviewerID int `json:"reviewer_id"`

	/* 多维度评分指标 */
	// 总体方面
	RankOverall int `json:"rank_overall"`

	// 内容、风格方面
	RankContent int `json:"rank_content"`

	// 工作量方面
	RankWorkload int `json:"rank_workload"`

	// 考核方面
	RankAssessment int `json:"rank_assessment"`

	// 点赞数
	UpvoteCount int `json:"upvote_count"`

	// 点踩数
	DownvoteCount int `json:"downvote_count"`
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
	ReviewID int `json:"review_id"`

	// 修改人 ID
	AlterBy int `json:"alter_by"`

	// 修改前的标题
	Title string `json:"title"`

	// 修改前的内容
	Content string `json:"content"`
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
