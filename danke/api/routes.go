package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/opentreehole/backend/danke/data"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api")
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", swagger.HandlerDefault)

	api := app.Group("/api")
	registerRoutes(api)
}

func registerRoutes(r fiber.Router) {
	r.Get("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to danke API"})
	})

	// v1
	// Course
	r.Get("/courses", ListCoursesV1)
	r.Get("/courses/:id<int>", GetCourseV1)
	r.Post("/courses", AddCourseV1)

	// CourseGroup
	r.Get("/group/:id<int>", GetCourseGroupV1)
	r.Get("/courses/hash", GetCourseGroupHashV1)
	r.Get("/courses/refresh", RefreshCourseGroupHashV1)

	// Review
	r.Get("/reviews/:id<int>", GetReviewV1)
	r.Get("/courses/:id<int>/reviews", ListReviewsV1)
	r.Post("/courses/:id<int>/reviews", CreateReviewV1)
	r.Put("/reviews/:id<int>", ModifyReviewV1)
	r.Patch("/reviews/:id<int>/_modify", ModifyReviewV1)
	r.Patch("/reviews/:id<int>", VoteForReviewV1)
	r.Get("/reviews/me", ListMyReviewsV1)
	r.Get("/reviews/random", GetRandomReviewV1)
	r.Delete("/reviews/:id<int>", DeleteReviewV1)

	// v3
	// CourseGroup
	r.Get("/v3/course_groups/search", SearchCourseGroupV3)
	r.Get("/v3/course_groups/:id<int>", GetCourseGroupV3)

	// static
	r.Get("/static/cedict_ts.u8", func(c *fiber.Ctx) error {
		return c.Send(data.CreditTs)
	})

	r.Get("/v3/reviews/_sensitive", ListSensitiveReviews)
	r.Put("/v3/reviews/:id<int>/_sensitive", ModifyReviewSensitive)
	r.Patch("/v3/reviews/:id<int>/_sensitive", ModifyReviewSensitive)
}
