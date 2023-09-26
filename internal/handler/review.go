package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/internal/service"
)

type ReviewHandler interface {
	RouteRegister
}

type reviewHandler struct {
	*Handler
	reviewService     service.ReviewService
	accountRepository repository.AccountRepository
}

func NewReviewHandler(
	handler *Handler,
	reviewService service.ReviewService,
	accountRepository repository.AccountRepository,
) ReviewHandler {
	return &reviewHandler{
		Handler:           handler,
		reviewService:     reviewService,
		accountRepository: accountRepository,
	}
}

func (h *reviewHandler) RegisterRoute(router fiber.Router) {
	router.Post("/courses/:course_id<int>/reviews", h.CreateReviewV1)
	router.Put("/reviews/:review_id<int>", h.ModifyReviewV1)
	router.Patch("/reviews/:review_id<int>", h.VoteForReviewV1)
	router.Get("/reviews/me", h.ListMyReviewsV1)
	router.Get("/reviews/random", h.GetRandomReviewV1)
}

// CreateReviewV1 godoc
// @Summary create a review
// @Description create a review
// @Tags Review
// @Accept json
// @Produce json
// @Param json body schema.CreateReviewV1Request true "json"
// @Param course_id path int true "course id"
// @Router /courses/{course_id}/reviews [post]
// @Success 200 {object} schema.ReviewV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
func (h *reviewHandler) CreateReviewV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return
	}

	var req schema.CreateReviewV1Request
	err = h.ValidateBody(c, &req)
	if err != nil {
		return
	}

	courseID, err := c.ParamsInt("course_id")
	if err != nil {
		return
	}

	response, err := h.reviewService.CreateReview(c.Context(), user, courseID, &req)
	if err != nil {
		return
	}

	return c.JSON(response)
}

// ModifyReviewV1 godoc
// @Summary modify a review
// @Description modify a review, admin or owner can modify
// @Tags Review
// @Accept json
// @Produce json
// @Param json body schema.ModifyReviewV1Request true "json"
// @Param review_id path int true "review id"
// @Router /reviews/{review_id} [put]
// @Success 200 {object} schema.ReviewV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
func (h *reviewHandler) ModifyReviewV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return
	}

	var req schema.ModifyReviewV1Request
	err = h.ValidateBody(c, &req)
	if err != nil {
		return
	}

	reviewID, err := c.ParamsInt("review_id")
	if err != nil {
		return
	}

	response, err := h.reviewService.ModifyReview(c.Context(), user, reviewID, &req)
	if err != nil {
		return
	}

	return c.JSON(response)
}

// VoteForReviewV1 godoc
// @Summary vote for a review
// @Description vote for a review
// @Tags Review
// @Accept json
// @Produce json
// @Param json body schema.VoteForReviewV1Request true "json"
// @Param review_id path int true "review id"
// @Router /reviews/{review_id} [patch]
// @Success 200 {object} schema.ReviewV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
func (h *reviewHandler) VoteForReviewV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return
	}

	var req schema.VoteForReviewV1Request
	err = h.ValidateBody(c, &req)
	if err != nil {
		return
	}

	reviewID, err := c.ParamsInt("review_id")
	if err != nil {
		return
	}

	response, err := h.reviewService.VoteForReviewV1(c.Context(), user, reviewID, &req)
	if err != nil {
		return
	}

	return c.JSON(response)
}

// ListMyReviewsV1 godoc
// @Summary list my reviews
// @Description list my reviews, old version. load history and achievements, no `is_me` field
// @Tags Review
// @Accept json
// @Produce json
// @Router /reviews/me [get]
// @Success 200 {array} schema.MyReviewV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
func (h *reviewHandler) ListMyReviewsV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return
	}

	response, err := h.reviewService.ListMyReviewsV1(c.Context(), user)
	if err != nil {
		return
	}

	return c.JSON(response)
}

// GetRandomReviewV1 godoc
// @Summary get random review
// @Description get random review
// @Tags Review
// @Accept json
// @Produce json
// @Router /reviews/random [get]
// @Success 200 {object} schema.RandomReviewV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
func (h *reviewHandler) GetRandomReviewV1(c *fiber.Ctx) (err error) {
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return
	}

	response, err := h.reviewService.GetRandomReviewV1(c.Context(), user)
	if err != nil {
		return
	}

	return c.JSON(response)
}
