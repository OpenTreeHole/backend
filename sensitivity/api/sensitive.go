package api

import (
	"github.com/gofiber/fiber/v2"
	. "github.com/opentreehole/backend/common"
	"github.com/opentreehole/backend/common/sensitive"
	. "github.com/opentreehole/backend/sensitivity/schema"
	"time"
)

func CheckSensitiveText(c *fiber.Ctx) (err error) {
	var request SensitiveCheckRequest
	err = ValidateBody(c, &request)
	if err != nil {
		return err
	}

	sensitiveResp, err := sensitive.CheckSensitive(sensitive.ParamsForCheck{
		Content:  request.Content,
		Id:       time.Now().UnixNano(),
		TypeName: sensitive.TypeTitle,
	})

	if err != nil {
		return nil
	}
	return c.JSON(SensitiveCheckResponse{ResponseForCheck: *sensitiveResp})
}
