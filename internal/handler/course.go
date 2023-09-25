package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/internal/repository"
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
}

// ListCoursesV1 godoc
// @Summary /courses
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
	ctx := context.WithValue(c.Context(), "FiberCtx", c)

	_, err = h.accountRepository.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	response, err := h.courseService.ListCoursesV1(ctx)
	if err != nil {
		return err
	}

	return c.JSON(response)
}
