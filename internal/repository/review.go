package repository

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
)

type ReviewRepository interface {
	Repository

	FindReviewsByCourseIDs(ctx context.Context, courseIDs []int, condition func(db *gorm.DB) *gorm.DB) (reviews []*model.Review, err error)
	FindReviews(ctx context.Context, condition func(db *gorm.DB) *gorm.DB) (reviews []*model.Review, err error)
	GetReviewByID(ctx context.Context, id int) (review *model.Review, err error)
	GetReview(ctx context.Context, condition func(tx *gorm.DB) *gorm.DB) (review *model.Review, err error)
	FindReviewVotes(ctx context.Context, reviewIDs []int, userIDs []int) (votes []*model.ReviewVote, err error)

	CreateReview(ctx context.Context, review *model.Review) (err error)

	UpdateReview(ctx context.Context, userID int, oldReview *model.Review, newReview *model.Review) (err error)
	UpdateReviewVote(ctx context.Context, userID int, review *model.Review, data int) (err error)
}

type reviewRepository struct {
	Repository
}

func NewReviewRepository(repository Repository) ReviewRepository {
	return &reviewRepository{Repository: repository}
}

/* 接口实现 */

func (r *reviewRepository) FindReviewsByCourseIDs(
	ctx context.Context,
	courseIDs []int,
	condition func(db *gorm.DB) *gorm.DB,
) (
	reviews []*model.Review,
	err error,
) {
	reviews = make([]*model.Review, 0, 5)
	err = condition(r.GetDB(ctx).Where("course_id IN ?", courseIDs)).Find(&reviews).Error
	return
}

func (r *reviewRepository) FindReviewVotes(ctx context.Context, reviewIDs []int, userIDs []int) (votes []*model.ReviewVote, err error) {
	votes = make([]*model.ReviewVote, 0, 5)
	if len(reviewIDs) == 0 && len(userIDs) == 0 {
		return
	}
	db := r.GetDB(ctx)
	if len(reviewIDs) > 0 {
		db = db.Where("review_id IN ?", reviewIDs)
	}
	if len(userIDs) > 0 {
		db = db.Where("user_id IN ?", userIDs)
	}
	err = db.Find(&votes).Error
	return
}

func (r *reviewRepository) FindReviews(ctx context.Context, condition func(db *gorm.DB) *gorm.DB) (reviews []*model.Review, err error) {
	reviews = make([]*model.Review, 0, 5)
	err = condition(r.GetDB(ctx)).Preload("History").
		Preload("UserAchievements.Achievement").Find(&reviews).Error
	return
}

func (r *reviewRepository) GetReview(ctx context.Context, condition func(tx *gorm.DB) *gorm.DB) (review *model.Review, err error) {
	review = new(model.Review)
	err = condition(r.GetDB(ctx)).Preload("History").
		Preload("UserAchievements.Achievement").First(review).Error
	return
}

func (r *reviewRepository) CreateReview(ctx context.Context, review *model.Review) (err error) {
	return r.Transaction(ctx, func(ctx context.Context) error {
		// create review
		err = r.GetDB(ctx).Create(review).Error
		if err != nil {
			return err
		}

		// update course review count
		err = r.GetDB(ctx).Model(&model.Course{ID: review.CourseID}).
			Update("review_count", gorm.Expr("review_count + 1")).Error
		if err != nil {
			return err
		}

		// update course_group review count
		if review.Course != nil {
			err = r.GetDB(ctx).Model(&model.CourseGroup{ID: review.Course.CourseGroupID}).
				Update("review_count", gorm.Expr("review_count + 1")).Error
			if err != nil {
				return err
			}
		} else {
			var reviewGroupID int
			err = r.GetDB(ctx).Model(&model.Course{}).Select("course_group_id").
				Where("id = ?", review.CourseID).Scan(&reviewGroupID).Error
			if err != nil {
				return err
			}

			err = r.GetDB(ctx).Model(&model.CourseGroup{ID: reviewGroupID}).
				Update("review_count", gorm.Expr("review_count + 1")).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *reviewRepository) GetReviewByID(ctx context.Context, id int) (review *model.Review, err error) {
	review = new(model.Review)
	err = r.GetDB(ctx).Preload("History").
		Preload("UserAchievements.Achievement").First(review, id).Error
	return
}

func (r *reviewRepository) UpdateReview(ctx context.Context, userID int, oldReview *model.Review, newReview *model.Review) (err error) {
	// 存储到 review_history 中
	err = r.GetDB(ctx).Create(&model.ReviewHistory{
		ReviewID: oldReview.ID,
		Title:    oldReview.Title,
		Content:  oldReview.Content,
		AlterBy:  userID,
	}).Error
	if err != nil {
		return
	}

	// 更新 review
	return r.GetDB(ctx).Model(oldReview).Updates(map[string]any{
		"title":           newReview.Title,
		"content":         newReview.Content,
		"modify_count":    gorm.Expr("modify_count + 1"),
		"rank_overall":    newReview.Rank.Overall,
		"rank_content":    newReview.Rank.Content,
		"rank_workload":   newReview.Rank.Workload,
		"rank_assessment": newReview.Rank.Assessment,
	}).Error
}

func (r *reviewRepository) UpdateReviewVote(ctx context.Context, userID int, review *model.Review, data int) (err error) {
	return r.Transaction(ctx, func(ctx context.Context) error {
		if data == 0 {
			err = r.GetDB(ctx).Where("review_id = ? AND user_id = ?", review.ID, userID).Delete(&model.ReviewVote{}).Error
		} else {
			err = r.GetDB(ctx).Save(&model.ReviewVote{
				UserID:   userID,
				ReviewID: review.ID,
				Data:     data,
			}).Error
		}
		if err != nil {
			return err
		}

		// update review vote count
		return r.GetDB(ctx).Exec(`
update review
set upvote_count = (select count(*) from review_vote where review_id = @review_id and data = 1), 
downvote_count = (select count(*) from review_vote where review_id = @review_id and data = -1) 
where id = @review_id`, sql.Named("review_id", review.ID)).Error
	})
}
