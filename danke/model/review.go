package model

import (
	"errors"
	"github.com/opentreehole/backend/common"
	"gorm.io/gorm"
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
	ID                  int                `json:"id"`
	CreatedAt           time.Time          `json:"created_at" gorm:"not null"`
	UpdatedAt           time.Time          `json:"updated_at" gorm:"not null"`
	CourseID            int                `json:"course_id" gorm:"not null;index"`
	Course              *Course            `json:"course"`
	Title               string             `json:"title" gorm:"not null"`
	Content             string             `json:"content" gorm:"not null"`
	ReviewerID          int                `json:"reviewer_id" gorm:"not null;index"`
	Rank                *ReviewRank        `json:"rank" gorm:"embedded;embeddedPrefix:rank_"`
	UpvoteCount         int                `json:"upvote_count" gorm:"not null;default:0"`
	DownvoteCount       int                `json:"downvote_count" gorm:"not null;default:0"`
	ModifyCount         int                `json:"modify_count" gorm:"not null;default:0"`
	History             ReviewHistoryList  `json:"-"`
	Vote                ReviewVoteList     `json:"-" gorm:"foreignKey:ReviewID;references:ID"`
	UserAchievements    []*UserAchievement `json:"-" gorm:"foreignKey:UserID;references:ReviewerID"`
	DeletedAt           gorm.DeletedAt     `json:"deleted_at" gorm:"index"`
	IsSensitive         bool               `json:"is_sensitive"`
	IsActuallySensitive *bool              `json:"is_actually_sensitive"`
	SensitiveDetail     string             `json:"sensitive_detail,omitempty"`
}

type FindReviewOption struct {
	PreloadHistory     bool
	PreloadAchievement bool
	PreloadVote        bool
	UserID             int
}

func (o FindReviewOption) setQuery(querySet *gorm.DB) *gorm.DB {
	if o.PreloadHistory {
		querySet = querySet.Preload("History")
	}
	if o.PreloadAchievement {
		querySet = querySet.Preload("UserAchievements.Achievement")
	}
	if o.PreloadVote {
		if o.UserID != 0 {
			querySet = querySet.Preload("Vote", "user_id = ?", o.UserID)
		} else {
			querySet = querySet.Preload("Vote")
		}
	}
	return querySet
}

func (r *Review) Sensitive() bool {
	if r == nil {
		return false
	}
	if r.IsActuallySensitive != nil {
		return *r.IsActuallySensitive
	}
	return r.IsSensitive
}

func FindReviewByID(tx *gorm.DB, reviewID int, options ...FindReviewOption) (review *Review, err error) {
	var option FindReviewOption
	if len(options) > 0 {
		option = options[0]
	}

	querySet := option.setQuery(tx)
	err = querySet.First(&review, reviewID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.NotFound("评论不存在")
	}
	return
}

func (r *Review) LoadVoteListByUserID(userID int) (err error) {
	return ReviewList{r}.LoadVoteListByUserID(userID)
}

func (r *Review) Create(tx *gorm.DB) (err error) {
	return tx.Transaction(func(tx *gorm.DB) (err error) {
		// 查找 course
		var course Course
		err = DB.First(&course, r.CourseID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.NotFound("课程不存在")
			}
			return err
		}

		// 检查是否已经评教过
		var count int64
		err = tx.Model(&Review{}).
			Where("course_id = ? AND reviewer_id = ?", r.CourseID, r.ReviewerID).
			Count(&count).Error
		if err != nil {
			return
		}
		if count > 0 {
			return common.BadRequest("已经评教过，请勿重复评教")
		}

		// 创建评教
		err = tx.Create(r).Error
		if err != nil {
			return
		}

		// 更新课程评教数量
		err = tx.Model(&Course{ID: r.CourseID}).
			Update("review_count", gorm.Expr("review_count + 1")).Error
		if err != nil {
			return
		}

		// 更新课程组评教数量
		if r.Course != nil {
			err = tx.Model(&CourseGroup{ID: r.Course.CourseGroupID}).
				Update("review_count", gorm.Expr("review_count + 1")).Error
			if err != nil {
				return err
			}
		} else {
			var reviewGroupID int
			err = tx.Model(&Course{}).Select("course_group_id").
				Where("id = ?", r.CourseID).Scan(&reviewGroupID).Error
			if err != nil {
				return err
			}

			err = tx.Model(&CourseGroup{ID: reviewGroupID}).
				Update("review_count", gorm.Expr("review_count + 1")).Error
			if err != nil {
				return err
			}
		}
		return
	})

}

func (r *Review) Update(tx *gorm.DB, newReview Review) (err error) {
	// 记录修改历史
	var history ReviewHistory
	history.FromReview(r)

	// 更新评教
	modified := false
	if newReview.Title != "" {
		r.Title = newReview.Title
		modified = true
	}
	if newReview.Content != "" {
		r.Content = newReview.Content
		modified = true
	}
	if newReview.Rank != nil {
		r.Rank = newReview.Rank
		modified = true
	}
	if !modified {
		return common.BadRequest("没有修改内容")
	}
	r.IsSensitive = newReview.IsSensitive
	r.IsActuallySensitive = newReview.IsActuallySensitive
	r.SensitiveDetail = newReview.SensitiveDetail

	r.ModifyCount++
	err = tx.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Model(&Review{ID: r.ID}).
			Select("Title", "Content", "Rank", "ModifyCount", "IsSensitive", "IsActuallySensitive", "SensitiveDetail").
			Updates(r).Error
		if err != nil {
			return
		}
		err = tx.Create(&history).Error
		return
	})
	return
}

type ReviewList []*Review

func (l ReviewList) LoadVoteListByUserID(userID int) (err error) {
	reviewIDs := make([]int, 0)
	for _, review := range l {
		reviewIDs = append(reviewIDs, review.ID)
	}
	var votes ReviewVoteList
	err = DB.
		Where("review_id IN ?", reviewIDs).
		Where("user_id = ?", userID).
		Find(&votes).Error
	if err != nil {
		return err
	}

	for _, vote := range votes {
		for _, review := range l {
			if review.ID == vote.ReviewID {
				review.Vote = append(review.Vote, vote)
			}
		}
	}
	return
}

// ReviewHistory 评教修改历史
type ReviewHistory struct {
	ID                  int       `json:"id"`
	CreatedAt           time.Time `json:"created_at"` // 创建时间，原本是 time_created
	UpdatedAt           time.Time `json:"updated_at"` // 更新时间，原本是 time_updated
	ReviewID            int       `json:"review_id" gorm:"not null;index"`
	AlterBy             int       `json:"alter_by" gorm:"not null"` // 修改人 ID
	Title               string    `json:"title" gorm:"not null"`    // 修改前的标题
	Content             string    `json:"content" gorm:"not null"`  // 修改前的内容
	IsSensitive         bool      `json:"is_sensitive"`
	IsActuallySensitive *bool     `json:"is_actual_sensitive"`
	SensitiveDetail     string    `json:"sensitive_detail,omitempty"`
}

func (h *ReviewHistory) FromReview(review *Review) {
	h.ReviewID = review.ID
	h.AlterBy = review.ReviewerID
	h.Title = review.Title
	h.Content = review.Content
	h.IsSensitive = review.IsSensitive
	h.IsActuallySensitive = review.IsActuallySensitive
	h.SensitiveDetail = review.SensitiveDetail
}

type ReviewHistoryList []*ReviewHistory

// ReviewVote 评教点赞/点踩详情
type ReviewVote struct {
	UserID   int `json:"user_id" gorm:"primaryKey"`   // 点赞或点踩人的 ID
	ReviewID int `json:"review_id" gorm:"primaryKey"` // 评教 ID
	Data     int `json:"data"`                        // 1 为点赞，-1 为点踩
}

type ReviewVoteList []*ReviewVote
