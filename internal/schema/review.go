package schema

import (
	"time"

	"github.com/jinzhu/copier"

	"github.com/opentreehole/backend/internal/model"
)

// ReviewRankV1 旧版本评分
type ReviewRankV1 struct {
	Overall    int `json:"overall" validate:"min=1,max=5"`    // 总体方面
	Content    int `json:"content" validate:"min=1,max=5"`    // 内容、风格方面
	Workload   int `json:"workload" validate:"min=1,max=5"`   // 工作量方面
	Assessment int `json:"assessment" validate:"min=1,max=5"` // 考核方面
}

func (r *ReviewRankV1) FromModel(rank *model.ReviewRank) *ReviewRankV1 {
	err := copier.Copy(r, rank)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *ReviewRankV1) ToModel() (rank *model.ReviewRank) {
	rank = new(model.ReviewRank)
	err := copier.Copy(rank, r)
	if err != nil {
		panic(err)
	}
	return
}

// AchievementV1Response 旧版本成就响应
type AchievementV1Response struct {
	Name       string    `json:"name"`        // 成就名称
	Domain     string    `json:"domain"`      // 成就域
	ObtainDate time.Time `json:"obtain_date"` // 获取日期
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
	ID          int                        `json:"id"`
	TimeCreated time.Time                  `json:"time_created" copier:"CreatedAt"` // 创建时间
	TimeUpdated time.Time                  `json:"time_updated" copier:"UpdatedAt"` // 更新时间
	Title       string                     `json:"title"`                           // 评教标题
	Content     string                     `json:"content"`                         // 评教内容
	ReviewerID  int                        `json:"reviewer_id"`                     // 评教者 ID
	Rank        *ReviewRankV1              `json:"rank"`                            // 评价
	Vote        int                        `json:"vote"`                            // 自己是否点赞或点踩，0 未操作，1 点赞，-1 点踩
	Remark      int                        `json:"remark"`                          // Remark = 点赞数 - 点踩数
	IsMe        bool                       `json:"is_me"`                           // 是否是自己的评教
	History     []*ReviewHistoryV1Response `json:"history"`                         // 修改历史，slices 必须非空
	Extra       UserExtraV1                `json:"extra"`                           // 额外信息
}

func (r *ReviewV1Response) FromModel(
	user *model.User,
	review *model.Review,
	votesMap map[int]*model.ReviewVote,
) *ReviewV1Response {
	err := copier.Copy(r, review)
	if err != nil {
		panic(err)
	}

	if user != nil {
		r.IsMe = user.ID == review.ReviewerID
	} else {
		r.IsMe = false
	}

	r.Rank = new(ReviewRankV1).FromModel(review.Rank)
	r.Remark = review.UpvoteCount - review.DownvoteCount
	if votesMap != nil && votesMap[review.ID] != nil {
		r.Vote = votesMap[review.ID].Data
	} else {
		r.Vote = 0
	}
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
	Title       string        `json:"title"`        // 旧标题
	Content     string        `json:"content"`      // 旧内容
	TimeCreated time.Time     `json:"time_created"` // 创建时间
	TimeUpdated time.Time     `json:"time_updated"` // 更新时间
	ReviewerID  int           `json:"reviewer_id"`  // 评教者
	Rank        *ReviewRankV1 `json:"rank"`         // 评价
	Remark      int           `json:"remark"`       // Remark = 点赞数 - 点踩数
}

// ReviewHistoryV1Response 旧版本评教修改历史响应
type ReviewHistoryV1Response struct {
	Time     time.Time        `json:"time"`     // 创建时间
	AlterBy  int              `json:"alter_by"` // 修改者
	Original *ReviewHistoryV1 `json:"original"` // 修改前的评教
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

type CreateReviewV1Request struct {
	Title   string       `json:"title" validate:"required,min=1,max=64"`
	Content string       `json:"content" validate:"required,min=1,max=10240"`
	Rank    ReviewRankV1 `json:"rank"`
}

type ModifyReviewV1Request = CreateReviewV1Request

func (r *CreateReviewV1Request) ToModel(reviewerID, courseID int) *model.Review {
	review := new(model.Review)
	err := copier.Copy(review, r)
	if err != nil {
		panic(err)
	}
	review.ReviewerID = reviewerID
	review.CourseID = courseID
	review.Rank = r.Rank.ToModel()
	return review
}

type VoteForReviewV1Request struct {
	Upvote bool `json:"upvote"`
}

type MyReviewV1Response struct {
	ID          int                        `json:"id"`
	Title       string                     `json:"title"`                           // 评教标题
	Content     string                     `json:"content"`                         // 评教内容
	History     []*ReviewHistoryV1Response `json:"history"`                         // 修改历史，slices 必须非空
	TimeCreated time.Time                  `json:"time_created" copier:"CreatedAt"` // 创建时间
	TimeUpdated time.Time                  `json:"time_updated" copier:"UpdatedAt"` // 更新时间
	ReviewerID  int                        `json:"reviewer_id"`                     // 评教者 ID
	Rank        *ReviewRankV1              `json:"rank"`                            // 评价
	Vote        int                        `json:"vote"`                            // 自己是否点赞或点踩，0 未操作，1 点赞，-1 点踩
	Remark      int                        `json:"remark"`                          // Remark = 点赞数 - 点踩数
	Extra       UserExtraV1                `json:"extra"`                           // 额外信息
	Course      *CourseV1Response          `json:"course,omitempty"`                // 课程信息
	GroupID     int                        `json:"group_id,omitempty"`              // 课程组 ID
}

func (r *MyReviewV1Response) FromModel(
	review *model.Review,
	votesMap map[int]*model.ReviewVote,
) *MyReviewV1Response {
	err := copier.Copy(r, review)
	if err != nil {
		panic(err)
	}

	r.Rank = new(ReviewRankV1).FromModel(review.Rank)
	r.Remark = review.UpvoteCount - review.DownvoteCount
	if votesMap != nil && votesMap[review.ID] != nil {
		r.Vote = votesMap[review.ID].Data
	} else {
		r.Vote = 0
	}
	r.History = make([]*ReviewHistoryV1Response, 0, len(review.History))
	for _, history := range review.History {
		r.History = append(r.History, new(ReviewHistoryV1Response).FromModel(review, history, r.Rank))
	}

	r.Extra.Achievements = make([]*AchievementV1Response, 0, len(review.UserAchievements))
	for _, userAchievement := range review.UserAchievements {
		r.Extra.Achievements = append(r.Extra.Achievements, new(AchievementV1Response).FromModel(userAchievement.Achievement, userAchievement))
	}

	// here course.Reviews is nil, so no need to send votesMap and user
	if review.Course != nil {
		r.Course = new(CourseV1Response).FromModel(nil, review.Course, nil)
		r.GroupID = review.Course.CourseGroupID
	}

	return r
}

type RandomReviewV1Response = MyReviewV1Response

/* V3 */

type ReviewRankV3 = ReviewRankV1

type AchievementV3Response = AchievementV1Response

type UserExtraV3 struct {
	Achievements []*AchievementV3Response `json:"achievements"`
}

type ReviewV3Response struct {
	ID            int           `json:"id"`
	CreatedAt     time.Time     `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time     `json:"updated_at"`     // 更新时间
	CourseID      int           `json:"course_id"`      // 课程 ID
	Title         string        `json:"title"`          // 评教标题
	Content       string        `json:"content"`        // 评教内容
	ReviewerID    int           `json:"reviewer_id"`    // 评教者 ID
	Rank          *ReviewRankV3 `json:"rank"`           // 评价
	MyVote        int           `json:"my_vote"`        // 自己是否点赞或点踩，0 未操作，1 点赞，-1 点踩
	UpvoteCount   int           `json:"upvote_count"`   // 点赞数
	DownvoteCount int           `json:"downvote_count"` // 点踩数
	IsMe          bool          `json:"is_me"`          // 是否是自己的评教
	Extra         UserExtraV3   `json:"extra"`          // 额外信息
}

func (r *ReviewV3Response) FromModel(
	user *model.User,
	review *model.Review,
	votesMap map[int]*model.ReviewVote,
) *ReviewV3Response {
	err := copier.Copy(r, review)
	if err != nil {
		panic(err)
	}

	if user != nil {
		r.IsMe = user.ID == review.ReviewerID
	} else {
		r.IsMe = false
	}

	r.Rank = new(ReviewRankV1).FromModel(review.Rank)
	if votesMap != nil && votesMap[review.ID] != nil {
		r.MyVote = votesMap[review.ID].Data
	} else {
		r.MyVote = 0
	}

	r.Extra.Achievements = make([]*AchievementV1Response, 0, len(review.UserAchievements))
	for _, userAchievement := range review.UserAchievements {
		if userAchievement.Achievement == nil {
			continue
		}
		r.Extra.Achievements = append(r.Extra.Achievements, new(AchievementV1Response).FromModel(userAchievement.Achievement, userAchievement))
	}
	return r

}
