package api

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router) {
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
	r.Post("/courses/:id<int>/reviews", CreateReviewV1)
	r.Put("/reviews/:id<int>", ModifyReviewV1)
	r.Patch("/reviews/:id<int>", VoteForReviewV1)
	r.Get("/reviews/me", ListMyReviewsV1)
	r.Get("/reviews/random", GetRandomReviewV1)

	// v3
	// CourseGroup
	r.Get("/v3/course_groups/search", SearchCourseGroupV3)
	r.Get("/v3/course_groups/:id<int>", GetCourseGroupV3)

	// static
	//router.Get("/static/cedict_ts.u8", func(c *fiber.Ctx) error {
	//	return c.Send(data.CreditTs)
	//})
}
