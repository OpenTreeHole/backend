package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/backend/common"
	. "github.com/opentreehole/backend/image_hosting/model"
	"log/slog"
	"strings"
)

func GetImage(c *fiber.Ctx) error {
	// to access the image in database
	// year := c.Params("year")
	// month := c.Params("month")
	// day := c.Params("day")
	// imageType := strings.Split(c.Params("imageIdentifier"), ".")[1]
	slog.LogAttrs(context.Background(), slog.LevelInfo, "getting image")
	var image ImageTable
	imageIdentifier := strings.Split(c.Params("identifier"), ".")[0]
	err := DB.First(&image, "image_identifier = ?", imageIdentifier)
	if err.Error != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Cannot find the image", slog.String("err", err.Error.Error()))
		return common.BadRequest("Cannot find the image")
	}
	slog.LogAttrs(context.Background(), slog.LevelInfo, "get image successfully", slog.String("image identifier", imageIdentifier))

	// browser will automatically transform the BLOB data to an image (no matter what extension)
	return c.Send(image.ImageFileData)
}
