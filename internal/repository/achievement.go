package repository

import (
	"context"

	"github.com/opentreehole/backend/internal/model"
)

type AchievementRepository interface {
	Repository

	FindAchievementsByUserID(ctx context.Context, userID int) (achievements []*model.Achievement, err error)
}

type achievementRepository struct {
	Repository
}

func NewAchievementRepository(repository Repository) AchievementRepository {
	return &achievementRepository{Repository: repository}
}

func (r *achievementRepository) FindAchievementsByUserID(ctx context.Context, userID int) (achievements []*model.Achievement, err error) {
	achievements = make([]*model.Achievement, 5)
	err = r.GetDB(ctx).Model(&achievements).
		Joins("JOIN user_achievements on achievement.id = user_achievements.achievement_id").
		Where("user_id = ?", userID).Find(&achievements).Error
	return
}
