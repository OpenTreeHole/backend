//! 迁移旧版本蛋壳到新版本
//! `review`.`history` -> `review_history`

package main

import (
	"time"

	"github.com/goccy/go-json"
	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/internal/model"
	"github.com/opentreehole/backend/internal/pkg/types"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/pkg/log"
)

type ReviewHistoryOld struct {
	Time     types.CustomTime `json:"time"`
	AlterBy  int              `json:"alter_by,omitempty"`
	Original any              `json:"original,omitempty"`
}

type ReviewRankOld struct {
	Overall    float64 `json:"overall"`
	Content    float64 `json:"content"`
	Workload   float64 `json:"workload"`
	Assessment float64 `json:"assessment"`
}

type ReviewHistoryOriginalOld struct {
	Title       string           `json:"title"`
	Content     string           `json:"content"`
	TimeCreated types.CustomTime `json:"time_created"`
	TimeUpdated types.CustomTime `json:"time_updated"`
	ReviewerID  int              `json:"reviewer_id"`
	Rank        *ReviewRankOld   `json:"rank"`
	Remark      int              `json:"remark"`
}

type ReviewOld struct {
	ID          int                 `json:"id"`
	TimeCreated time.Time           `json:"time_created"`
	TimeUpdated time.Time           `json:"time_updated"`
	Title       string              `json:"title"`
	Content     string              `json:"content"`
	ReviewerID  int                 `json:"reviewer_id"`
	Rank        *ReviewRankOld      `json:"rank" gorm:"serializer:json"`
	History     []*ReviewHistoryOld `json:"history" gorm:"serializer:json"`
	Upvoters    []int               `json:"upvoters" gorm:"serializer:json"`
	Downvoters  []int               `json:"downvoters" gorm:"serializer:json"`
}

func main() {
	conf := config.NewConfig()
	logger, cancel := log.NewLogger(conf)
	defer cancel()
	db := repository.NewDB(conf, logger)

	var (
		reviews []*ReviewOld
		err     error
	)

	err = db.Transaction(func(tx *gorm.DB) error {
		err = db.Table("review").FindInBatches(&reviews, 1000, func(tx *gorm.DB, batch int) error {
			// update History
			var (
				newHistory    []*model.ReviewHistory
				newReviewVote []*model.ReviewVote
			)
			for _, review := range reviews {
				// history
				for _, history := range review.History {
					var data []byte
					switch h := history.Original.(type) {
					case []byte:
						data = h
					case string:
						data = []byte(h)
					default:
						data, err = json.Marshal(h)
						if err != nil {
							return err
						}
					}
					var original ReviewHistoryOriginalOld
					err = json.Unmarshal(data, &original)
					if err != nil {
						return err
					}
					history.Original = original

					newHistory = append(newHistory, &model.ReviewHistory{
						CreatedAt: history.Time.Time,
						UpdatedAt: history.Time.Time,
						ReviewID:  review.ID,
						AlterBy:   history.AlterBy,
						Title:     original.Title,
						Content:   original.Content,
					})
				}

				// vote
				type ReviewVote struct {
					UserID   int
					ReviewID int
				}
				var reviewVoteMap = make(map[ReviewVote]int)
				for _, userID := range review.Downvoters {
					reviewVoteMap[ReviewVote{userID, review.ID}] = -1
				}
				for _, userID := range review.Upvoters {
					reviewVoteMap[ReviewVote{userID, review.ID}] = 1
				}

				for k, v := range reviewVoteMap {
					newReviewVote = append(newReviewVote, &model.ReviewVote{UserID: k.UserID, ReviewID: k.ReviewID, Data: v})
				}
			}

			err = tx.Create(&newHistory).Error
			if err != nil {
				return err
			}

			err = tx.Create(&newReviewVote).Error
			if err != nil {
				return err
			}

			return nil
		}).Error

		if err != nil {
			return err
		}

		// update course.review_count
		err = db.Exec(`update course set review_count = (select count(*) from review where review.course_id = course.id) where true`).Error
		if err != nil {
			return err
		}

		// update review.upvote_count and review.downvote_count
		// extract review.rank_* from review.rank
		err = db.Exec(`update review 
set upvote_count = (select count(*) from review_vote where review_vote.review_id = review.id and review_vote.data = 1), 
    downvote_count = (select count(*) from review_vote where review_vote.review_id = review.id and review_vote.data = -1),
	rank_overall = JSON_EXTRACT(review.rank, '$.overall'),
	rank_content = JSON_EXTRACT(review.rank, '$.content'),
	rank_assessment = JSON_EXTRACT(review.rank, '$.assessment'),
	rank_workload = JSON_EXTRACT(review.rank, '$.workload')
where true`).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
