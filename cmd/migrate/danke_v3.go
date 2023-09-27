package migrate

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

func DankeV3() {
	conf := config.NewConfig()
	logger, cancel := log.NewLogger(conf)
	defer cancel()
	db := repository.NewDB(conf, logger)

	var (
		reviews []*ReviewOld
		err     error
	)
	logger.Info("migrate danke v3")

	err = db.Transaction(func(tx *gorm.DB) error {
		m := tx.Migrator()

		if m.HasConstraint(&model.Course{}, "course_ibfk_1") {
			err = m.DropConstraint(&model.Course{}, "course_ibfk_1")
			if err != nil {
				return err
			}
		}

		if m.HasConstraint(&model.Review{}, "review_ibfk_1") {
			err = m.DropConstraint(&model.Review{}, "review_ibfk_1")
			if err != nil {
				return err
			}
		}

		if m.HasTable("coursegroup") {
			err = m.RenameTable("coursegroup", "course_group")
			if err != nil {
				return err
			}
		}

		if m.HasColumn(&model.Course{}, "coursegroup_id") {
			err = m.RenameColumn(&model.Course{}, "coursegroup_id", "course_group_id")
			if err != nil {
				return err
			}
		}

		if m.HasColumn(&model.Review{}, "time_created") {
			err = m.RenameColumn(&model.Review{}, "time_created", "created_at")
			if err != nil {
				return err
			}
		}

		if m.HasColumn(&model.Review{}, "time_updated") {
			err = m.RenameColumn(&model.Review{}, "time_updated", "updated_at")
			if err != nil {
				return err
			}
		}

		if m.HasIndex(&model.Course{}, "coursegroup_id") {
			err = m.RenameIndex(&model.Course{}, "coursegroup_id", "idx_course_course_group_id")
			if err != nil {
				return err
			}
		}

		if m.HasIndex(&model.Review{}, "course_id") {
			err = m.RenameIndex(&model.Review{}, "course_id", "idx_review_course_id")
			if err != nil {
				return err
			}
		}

		if m.HasIndex(&model.Review{}, "reviewer_id") {
			err = m.RenameIndex(&model.Review{}, "reviewer_id", "idx_review_reviewer_id")
			if err != nil {
				return err
			}
		}

		err = m.AutoMigrate(
			&model.Course{},
			&model.Review{},
			&model.ReviewHistory{},
			&model.ReviewVote{},
			&model.CourseGroup{},
			&model.UserAchievement{},
		)
		if err != nil {
			return err
		}

		logger.Info("migrate danke v3: update review history and vote")
		err = tx.Table("review").FindInBatches(&reviews, 1000, func(tx *gorm.DB, batch int) error {
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

			if len(newHistory) > 0 {
				err = tx.Create(&newHistory).Error
				if err != nil {
					return err
				}
			}

			if len(newReviewVote) > 0 {
				err = tx.Create(&newReviewVote).Error
				if err != nil {
					return err
				}
			}

			return nil
		}).Error

		if err != nil {
			return err
		}

		// update course.review_count
		err = tx.Exec(`update course set review_count = (select count(*) from review where review.course_id = course.id) where true`).Error
		if err != nil {
			return err
		}

		// update review.upvote_count and review.downvote_count
		// extract review.rank_* from review.rank
		err = tx.Exec(`update review 
set upvote_count = (select count(*) from review_vote where review_vote.review_id = review.id and review_vote.data = 1), 
    downvote_count = (select count(*) from review_vote where review_vote.review_id = review.id and review_vote.data = -1),
	rank_overall = JSON_EXTRACT(review.rank, '$.overall'),
	rank_content = JSON_EXTRACT(review.rank, '$.content'),
	rank_assessment = JSON_EXTRACT(review.rank, '$.assessment'),
	rank_workload = JSON_EXTRACT(review.rank, '$.workload'),
    modify_count = (select count(*) from review_history where review_history.review_id = review.id)
where true`).Error
		if err != nil {
			return err
		}

		if m.HasColumn(&model.Review{}, "upvoters") {
			err = m.DropColumn(&model.Review{}, "upvoters")
			if err != nil {
				return err
			}
		}

		if m.HasColumn(&model.Review{}, "downvoters") {
			err = m.DropColumn(&model.Review{}, "downvoters")
			if err != nil {
				return err
			}
		}

		if m.HasColumn(&model.Review{}, "rank") {
			err = m.DropColumn(&model.Review{}, "rank")
			if err != nil {
				return err
			}
		}

		if m.HasColumn(&model.Review{}, "history") {
			err = m.DropColumn(&model.Review{}, "history")
			if err != nil {
				return err
			}
		}

		err = m.DropTable("userextra", "aerich", "seaql_migrations")
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
