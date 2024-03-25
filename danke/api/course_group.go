package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	. "github.com/opentreehole/backend/common"
	. "github.com/opentreehole/backend/danke/model"
	. "github.com/opentreehole/backend/danke/schema"
	"gorm.io/gorm"
)

// GetCourseGroupV1 godoc
// @Summary /group/{group_id}
// @Description get a course group, old version or v1 version
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Deprecated
// @Router /group/{id} [get]
// @Param id path int true "course group id"
// @Success 200 {object} schema.CourseGroupV1Response
// @Failure 400 {object} common.HttpError
// @Failure 404 {object} common.HttpBaseError
// @Failure 500 {object} common.HttpBaseError
func GetCourseGroupV1(c *fiber.Ctx) (err error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return err
	}

	groupID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// 获取课程组，同时加载课程
	// 这里不预加载课程的评论，因为评论作为动态的数据，应该独立作缓存，提高缓存粒度和缓存更新频率
	var courseGroup CourseGroup
	err = DB.Preload("Courses").First(&courseGroup, groupID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFound("课程组不存在")
		}
		return err
	}

	// 获取课程组的所有课程的所有评论，同时加载评论的历史记录和用户成就
	err = courseGroup.Courses.LoadReviewList(DB, true, true)
	if err != nil {
		return err
	}

	// 获取课程组的所有课程的所有评论的自己的投票
	err = courseGroup.Courses.AllReviewList().LoadVoteListByUserID(user.ID)
	if err != nil {
		return err
	}

	return c.JSON(new(CourseGroupV1Response).FromModel(user, &courseGroup))
}

// GetCourseGroupHashV1 godoc
// @Summary get course group hash
// @Description get course group hash
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Router /courses/hash [get]
// @Success 200 {object} schema.CourseGroupHashV1Response
// @Failure 400 {object} common.HttpError
// @Failure 404 {object} common.HttpBaseError
// @Failure 500 {object} common.HttpBaseError
func GetCourseGroupHashV1(c *fiber.Ctx) (err error) {
	_, err = GetCurrentUser(c)
	if err != nil {
		return err
	}

	// 获取课程组哈希
	_, hash, err := FindGroupsWithCourses(false)
	if err != nil {
		return
	}

	return c.JSON(CourseGroupHashV1Response{Hash: hash})
}

// RefreshCourseGroupHashV1 godoc
// @Summary refresh course group hash
// @Description refresh course group hash, admin only
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Router /courses/refresh [get]
// @Success 418
// @Failure 400 {object} common.HttpError
// @Failure 404 {object} common.HttpBaseError
// @Failure 500 {object} common.HttpBaseError
func RefreshCourseGroupHashV1(c *fiber.Ctx) (err error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return Forbidden()
	}

	// 刷新课程组哈希
	_, _, err = FindGroupsWithCourses(true)
	if err != nil {
		return
	}

	return c.SendStatus(fiber.StatusTeapot)
}

// SearchCourseGroupV3 godoc
// @Summary search course group
// @Description search course group, no courses
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Router /v3/course_groups/search [get]
// @Param request query schema.CourseGroupSearchV3Request true "search query"
// @Success 200 {object} common.PagedResponse[schema.CourseGroupV3Response, any]
// @Failure 400 {object} common.HttpError
// @Failure 404 {object} common.HttpBaseError
// @Failure 500 {object} common.HttpBaseError
func SearchCourseGroupV3(c *fiber.Ctx) (err error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return err
	}

	var request CourseGroupSearchV3Request
	err = ValidateQuery(c, &request)
	if err != nil {
		return
	}

	var (
		page     = request.Page
		pageSize = request.PageSize
		query    = request.Query
	)

	var querySet = DB
	if CourseCodeRegexp.MatchString(query) {
		querySet = querySet.Where("code LIKE ?", query+"%")
	} else {
		querySet = querySet.Where("name LIKE ?", "%"+query+"%")
	}
	if page > 0 {
		if pageSize == 0 {
			pageSize = 10
		}
		querySet = querySet.Limit(pageSize).Offset((page - 1) * pageSize)
	} else {
		page = 1
		if pageSize > 0 {
			querySet = querySet.Limit(pageSize)
		}
	}
	querySet = querySet.Order("id")

	var courseGroups CourseGroupList
	err = querySet.Find(&courseGroups).Error
	if err != nil {
		return err
	}

	items := make([]*CourseGroupV3Response, 0, len(courseGroups))
	for _, group := range courseGroups {
		items = append(items, new(CourseGroupV3Response).FromModel(user, group))
	}

	return c.JSON(PagedResponse[CourseGroupV3Response, any]{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetCourseGroupV3 godoc
// @Summary get a course group
// @Description get a course group, v3 version
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Router /v3/course_groups/{id} [get]
// @Param id path int true "course group id"
// @Success 200 {object} schema.CourseGroupV3Response
// @Failure 400 {object} common.HttpError
// @Failure 404 {object} common.HttpBaseError
// @Failure 500 {object} common.HttpBaseError
func GetCourseGroupV3(c *fiber.Ctx) (err error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return err
	}

	groupID, err := c.ParamsInt("id")
	if err != nil {
		return
	}

	var courseGroup CourseGroup
	err = DB.Preload("Courses").First(&courseGroup, groupID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFound("课程组不存在")
		}
		return err
	}

	// 获取课程组的所有课程的所有评论，同时加载评论的历史记录和用户成就
	err = courseGroup.Courses.LoadReviewList(DB, true, true)
	if err != nil {
		return err
	}

	// 获取课程组的所有课程的所有评论的自己的投票
	err = courseGroup.Courses.AllReviewList().LoadVoteListByUserID(user.ID)
	if err != nil {
		return err
	}

	return c.JSON(new(CourseGroupV3Response).FromModel(user, &courseGroup))
}
