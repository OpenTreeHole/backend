package schema

import (
	"time"

	"github.com/jinzhu/copier"

	"github.com/opentreehole/backend/internal/model"
)

// ReviewRankV1 旧版本评分
type ReviewRankV1 struct {
	// 总体方面
	Overall int `json:"overall"`

	// 内容、风格方面
	Content int `json:"content"`

	// 工作量方面
	Workload int `json:"workload"`

	// 考核方面
	Assessment int `json:"assessment"`
}

func (r *ReviewRankV1) FromModel(rank *model.ReviewRank) *ReviewRankV1 {
	err := copier.Copy(r, rank)
	if err != nil {
		panic(err)
	}
	return r
}

// AchievementV1Response 旧版本成就响应
type AchievementV1Response struct {
	// 成就名称
	Name string `json:"name"`

	// 成就域
	Domain string `json:"domain"`

	// 获取日期
	ObtainDate time.Time `json:"obtain_date"`
}

func (r *AchievementV1Response) FromModel(
	achievement *model.Achievement,
	userAchievement *model.UserAchievement,
) *AchievementV1Response {
	err := copier.Copy(r, userAchievement)
	if err != nil {
		panic(err)
	}

	r.Name = achievement.Name
	r.Domain = achievement.Domain

	return r
}

// UserExtraV1 旧版本用户额外信息
type UserExtraV1 struct {
	// 用户成就，slices 必须非空
	Achievements []*AchievementV1Response `json:"achievements"`
}

// ReviewV1Response 旧版本评教响应
type ReviewV1Response struct {
	// 评教 ID
	ID int `json:"id"`

	// 创建时间
	TimeCreated time.Time `json:"time_created"`

	// 更新时间
	TimeUpdated time.Time `json:"time_updated"`

	// 评教标题
	Title string `json:"title"`

	// 评教内容
	Content string `json:"content"`

	// 评教者 ID
	ReviewerID int `json:"reviewer_id"`

	// 评价
	Rank *ReviewRankV1 `json:"rank"`

	// 自己是否点赞或点踩，0 未操作，1 点赞，-1 点踩
	Vote int `json:"vote"`

	// Remark = 点赞数 - 点踩数
	Remark int `json:"remark"`

	// 是否是自己的评教
	IsMe bool `json:"is_me"`

	// 修改历史，slices 必须非空
	History []*ReviewHistoryV1Response `json:"history"`

	// 额外信息
	Extra UserExtraV1 `json:"extra"`
}

func (r *ReviewV1Response) FromModel(
	user *model.User,
	review *model.Review,
) *ReviewV1Response {
	err := copier.Copy(r, review)
	if err != nil {
		panic(err)
	}

	r.IsMe = user.ID == review.ReviewerID
	r.Rank = new(ReviewRankV1).FromModel(review.Rank)
	// TODO: vote
	r.History = make([]*ReviewHistoryV1Response, 0, len(review.History))
	for _, history := range review.History {
		r.History = append(r.History, new(ReviewHistoryV1Response).FromModel(review, history, r.Rank))
	}

	r.Extra.Achievements = make([]*AchievementV1Response, 0, len(review.UserAchievements))
	for _, userAchievement := range review.UserAchievements {
		r.Extra.Achievements = append(r.Extra.Achievements, new(AchievementV1Response).FromModel(userAchievement.Achievement, userAchievement))
	}
	return r
}

type ReviewHistoryV1 struct {
	// 旧标题
	Title string `json:"title"`

	// 旧内容
	Content string `json:"content"`

	// 创建时间
	TimeCreated time.Time `json:"time_created"`

	// 更新时间
	TimeUpdated time.Time `json:"time_updated"`

	// 评教者
	ReviewerID int `json:"reviewer_id"`

	// 评价
	Rank *ReviewRankV1 `json:"rank"`

	// Remark = 点赞数 - 点踩数
	Remark int `json:"remark"`
}

// ReviewHistoryV1Response 旧版本评教修改历史响应
type ReviewHistoryV1Response struct {
	// 创建时间
	Time time.Time `json:"time"`

	// 修改者
	AlterBy int `json:"alter_by"`

	// 修改前的评教
	Original *ReviewHistoryV1 `json:"original"`
}

func (r *ReviewHistoryV1Response) FromModel(
	review *model.Review,
	history *model.ReviewHistory,
	rank *ReviewRankV1,
) *ReviewHistoryV1Response {
	r.Time = history.CreatedAt
	r.AlterBy = history.AlterBy
	r.Original = &ReviewHistoryV1{
		Title:       history.Title,
		Content:     history.Content,
		TimeCreated: review.CreatedAt,
		TimeUpdated: history.CreatedAt,
		ReviewerID:  review.ReviewerID,
		Rank:        rank,
		Remark:      review.UpvoteCount - review.DownvoteCount,
	}

	return r
}
