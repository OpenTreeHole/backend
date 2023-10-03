package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/data"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/internal/service"
)

type CourseGroupHandler interface {
	RouteRegister
}

type courseGroupHandler struct {
	*Handler
	courseGroupService service.CourseGroupService
	accountRepository  repository.AccountRepository
}

func NewCourseGroupHandler(
	handler *Handler,
	groupService service.CourseGroupService,
	accountRepository repository.AccountRepository,
) CourseGroupHandler {
	return &courseGroupHandler{
		Handler:            handler,
		courseGroupService: groupService,
		accountRepository:  accountRepository,
	}
}

func (h *courseGroupHandler) RegisterRoute(router fiber.Router) {
	router.Get("/group/:id<int>", h.GetCourseGroupV1)
	router.Get("/courses/hash", h.GetCourseGroupHashV1)
	router.Get("/courses/refresh", h.RefreshCourseGroupHashV1)

	// v3
	router.Get("/v3/course_groups/search", h.SearchCourseGroupV3)
	router.Get("/v3/course_groups/:id<int>", h.GetCourseGroupV3)

	// static
	router.Get("/static/cedict_ts.u8", func(c *fiber.Ctx) error {
		return c.Send(data.CreditTs)
	})
}

// GetCourseGroupV1 godoc
// @Summary /group/{group_id}
// @Description get a course group, old version or v1 version
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Deprecated
// @Router /group/{id} [get]
// @Param id path string true "course group id"
// @Success 200 {object} schema.CourseGroupV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
// @Failure 500 {object} schema.HttpBaseError
func (h *courseGroupHandler) GetCourseGroupV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return err
	}

	groupID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	response, err := h.courseGroupService.GetGroupByIDV1(c.Context(), user, groupID)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// GetCourseGroupHashV1 godoc
// @Summary get course group hash
// @Description get course group hash
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Router /courses/hash [get]
// @Success 200 {object} schema.CourseGroupHashV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
// @Failure 500 {object} schema.HttpBaseError
func (h *courseGroupHandler) GetCourseGroupHashV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	_, err = h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return
	}

	response, err := h.courseGroupService.GetCourseGroupHash(c.Context())
	if err != nil {
		return
	}

	return c.JSON(response)
}

// RefreshCourseGroupHashV1 godoc
// @Summary refresh course group hash
// @Description refresh course group hash, admin only
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Router /courses/refresh [get]
// @Success 418
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
// @Failure 500 {object} schema.HttpBaseError
func (h *courseGroupHandler) RefreshCourseGroupHashV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return schema.Forbidden()
	}

	err = h.courseGroupService.RefreshCourseGroupHash(c.Context())
	if err != nil {
		return err
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
// @Success 200 {object} schema.PagedResponse[schema.CourseGroupV3Response, any]
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
// @Failure 500 {object} schema.HttpBaseError
func (h *courseGroupHandler) SearchCourseGroupV3(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return
	}

	request := new(schema.CourseGroupSearchV3Request)
	err = h.ValidateQuery(c, request)
	if err != nil {
		return
	}

	response, err := h.courseGroupService.SearchCourseGroupV3(c.Context(), user, request)
	if err != nil {
		return
	}

	return c.JSON(response)
}

// GetCourseGroupV3 godoc
// @Summary /v3/course_groups/{group_id}
// @Description get a course group, v3 version
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Router /v3/course_groups/{id} [get]
// @Param id path string true "course group id"
// @Success 200 {object} schema.CourseGroupV3Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
// @Failure 500 {object} schema.HttpBaseError
func (h *courseGroupHandler) GetCourseGroupV3(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return
	}

	groupID, err := c.ParamsInt("id")
	if err != nil {
		return
	}

	response, err := h.courseGroupService.GetGroupByIDV3(c.Context(), user, groupID)
	if err != nil {
		return
	}

	return c.JSON(response)
}
