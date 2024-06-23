package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	. "github.com/opentreehole/backend/common"
	. "github.com/opentreehole/backend/danke/model"
	. "github.com/opentreehole/backend/danke/schema"
	"gorm.io/gorm"
)

// ListCoursesV1 godoc
// @Summary list courses
// @Description list all course_groups and courses, no reviews, old version or v1 version
// @Tags Course
// @Accept json
// @Produce json
// @Deprecated
// @Router /courses [get]
// @Success 200 {array} schema.CourseGroupV1Response
// @Failure 400 {object} common.HttpError
// @Failure 404 {object} common.HttpBaseError
func ListCoursesV1(c *fiber.Ctx) (err error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return err
	}

	groups, _, err := FindGroupsWithCourses(false)
	if err != nil {
		return err
	}

	response := make([]*CourseGroupV1Response, 0, len(groups))
	for _, group := range groups {
		response = append(response, new(CourseGroupV1Response).FromModel(user, group))
	}

	return c.JSON(response)
}

// GetCourseV1 godoc
// @Summary get a course
// @Description get a course with reviews, v1 version
// @Tags Course
// @Accept json
// @Produce json
// @Deprecated
// @Router /courses/{id} [get]
// @Param id path int true "course_id"
// @Success 200 {object} schema.CourseV1Response
// @Failure 400 {object} common.HttpError
// @Failure 404 {object} common.HttpBaseError
func GetCourseV1(c *fiber.Ctx) (err error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return
	}

	// 获取课程，课程的评论，评论的历史记录和用户成就
	var course Course
	err = DB.
		Preload("Reviews.History").
		Preload("Reviews.UserAchievements.Achievement").
		First(&course, id).Error
	if err != nil {
		return
	}

	// 获取课程的评论的自己的投票
	err = course.Reviews.LoadVoteListByUserID(user.ID)
	if err != nil {
		return err
	}

	return c.JSON(new(CourseV1Response).FromModel(user, &course))
}

// AddCourseV1 godoc
// @Summary add a course
// @Description add a course, admin only
// @Tags Course
// @Accept json
// @Produce json
// @Router /courses [post]
// @Param json body schema.CreateCourseV1Request true "json"
// @Success 200 {object} schema.CourseV1Response
// @Failure 400 {object} common.HttpError
// @Failure 500 {object} common.HttpBaseError
func AddCourseV1(c *fiber.Ctx) (err error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return err
	}

	var request CreateCourseV1Request
	err = ValidateBody(c, &request)
	if err != nil {
		return err
	}

	// 查找课程
	var course *Course
	err = DB.First(&course, "code_id = ?", request.CodeID).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		return BadRequest("该课程已存在")
	}

	// 根据 Code 查找课程组
	var courseGroup *CourseGroup
	err = DB.Preload("Courses").First(&courseGroup, "code = ?", request.Code).Error
	if err != nil {
		// 如果没有找到课程组，创建一个新的课程组
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		courseGroup = request.ToCourseGroupModel()
		err = DB.Create(&courseGroup).Error
		if err != nil {
			return err
		}
	}

	course = request.ToModel(courseGroup.ID)
	course.CourseGroup = courseGroup
	err = course.Create()
	if err != nil {
		return err
	}

	return c.JSON(new(CourseV1Response).FromModel(user, course))
}
