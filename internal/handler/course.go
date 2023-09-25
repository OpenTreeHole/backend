package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/internal/service"
)

type CourseHandler interface {
	RouteRegister
}

type courseHandler struct {
	*Handler
	courseService      service.CourseService
	courseGroupService service.CourseGroupService
	accountRepository  repository.AccountRepository
}

func NewCourseHandler(
	handler *Handler,
	courseService service.CourseService,
	courseGroupService service.CourseGroupService,
	accountRepository repository.AccountRepository,
) CourseHandler {
	return &courseHandler{
		Handler:            handler,
		courseService:      courseService,
		courseGroupService: courseGroupService,
		accountRepository:  accountRepository,
	}
}

func (h *courseHandler) RegisterRoute(router fiber.Router) {
	router.Get("/courses", h.ListCoursesV1)
	router.Get("/courses/:id", h.GetCourseV1)
	router.Post("/courses", h.CreateCourseV1)
}

// ListCoursesV1 godoc
// @Summary list courses
// @Description list all course_groups and courses, no reviews, old version or v1 version
// @Tags Course
// @Accept json
// @Produce json
// @Deprecated
// @Router /courses [get]
// @Success 200 {array} schema.CourseGroupV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
func (h *courseHandler) ListCoursesV1(c *fiber.Ctx) (err error) {
	//ctx := context.WithValue(c.Context(), "FiberCtx", c)
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return err
	}
	fmt.Printf("%+v", *user)

	response, err := h.courseService.ListCoursesV1(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// GetCourseV1 godoc
// @Summary get a course
// @Description get a course with reviews, old version or v1 version
// @Tags Course
// @Accept json
// @Produce json
// @Deprecated
// @Router /courses/{id} [get]
// @Param id path int true "course_id"
// @Success 200 {object} schema.CourseV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
func (h *courseHandler) GetCourseV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	response, err := h.courseService.GetCourseV1(c.Context(), user, id)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// CreateCourseV1 godoc
// @Summary create a course
// @Description create a course, only admin can create
// @Tags Course
// @Accept json
// @Produce json
// @Router /courses [post]
// @Param json body schema.CreateCourseV1Request true "json"
// @Success 200 {object} schema.CourseV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 500 {object} schema.HttpBaseError
func (h *courseHandler) CreateCourseV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return schema.Forbidden()
	}

	var request schema.CreateCourseV1Request
	err = h.ValidateBody(c, &request)
	if err != nil {
		return err
	}

	response, err := h.courseService.AddCourseV1(c.Context(), &request)
	if err != nil {
		return err
	}

	return c.JSON(response)
}
