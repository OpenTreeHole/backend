package api

import (
	"bytes"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/backend/common"
	. "github.com/opentreehole/backend/image_hosting/config"
	. "github.com/opentreehole/backend/image_hosting/model"
	. "github.com/opentreehole/backend/image_hosting/schema"
	. "github.com/opentreehole/backend/image_hosting/utils"
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"time"
)

func UploadImage(c *fiber.Ctx) error {
	slog.LogAttrs(context.Background(), slog.LevelInfo, "uploading image")

	// response to frontend
	var response CheveretoUploadResponse
	// the file uploaded by user  in the request body with the form-data key "source"
	file, err := c.FormFile("source")
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "No file uploaded", slog.String("err", err.Error()))
		return common.BadRequest("No file uploaded")
	}

	fileSize := file.Size
	maxSize := 10 * 1024 * 1024 // file should <= 10MB
	if int(fileSize) > maxSize {
		slog.LogAttrs(context.Background(), slog.LevelError, "File size is too large")
		return common.BadRequest("File size is too large")
	}

	fileExtension := strings.TrimPrefix(filepath.Ext(file.Filename), ".")

	if !IsAllowedExtension(fileExtension) {
		slog.LogAttrs(context.Background(), slog.LevelError, "File type not allowed.")
		return common.BadRequest("File type not allowed.")
	}

	fileContent, err := file.Open()
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "The uploaded file has some problems", slog.String("err", err.Error()))
		return common.BadRequest("The uploaded file has some problems")
	}

	imageData, err := io.ReadAll(fileContent)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "The uploaded file has some problems", slog.String("err", err.Error()))
		return common.BadRequest("The uploaded file has some problems")
	}

	imageIdentifier, err := GenerateIdentifier()
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Cannot generate image identifier", slog.String("err", err.Error()))
		return common.InternalServerError("Cannot generate image identifier")
	}

	var image ImageTable
	originalFileName := file.Filename
	result := DB.First(&image, "original_file_name = ?", originalFileName)
	if result.Error == nil {
		if bytes.Equal(image.ImageFileData, imageData) && strings.EqualFold(image.ImageType, fileExtension) {
			slog.LogAttrs(context.Background(), slog.LevelInfo, "The file has been uploaded before")
			imageIdentifier = image.ImageIdentifier
			imageUrl := viper.GetString(EnvHostName) + "/api/i/" + image.CreatedAt.Format("2006/01/02/") + imageIdentifier + "." + fileExtension
			response.StatusCode = 200
			response.StatusTxt = "The image has been uploaded before"
			response.Image = CheveretoImageInfo{
				Name:       imageIdentifier,
				Extension:  image.ImageType,
				Filename:   imageIdentifier + "." + fileExtension,
				Url:        imageUrl,
				DisplayUrl: imageUrl,
				Mime:       "image/" + fileExtension,
			}
			return c.JSON(&response)
		}
	}

	imageUrl := viper.GetString(EnvHostName) + "/api/i/" + time.Now().Format("2006/01/02/") + imageIdentifier + "." + fileExtension
	uploadedImage := &ImageTable{
		ImageIdentifier:  imageIdentifier,
		OriginalFileName: originalFileName,
		ImageType:        fileExtension,
		ImageFileData:    imageData,
	}
	err = DB.Create(&uploadedImage).Error

	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Database cannot store the image", slog.String("err", err.Error()))
		return common.InternalServerError("Database cannot store the image")
	}

	// if nothing went wrong
	response.StatusCode = 200
	response.StatusTxt = "Upload Success"
	response.Image = CheveretoImageInfo{
		Name:       imageIdentifier,
		Extension:  fileExtension,
		Filename:   imageIdentifier + "." + fileExtension,
		Url:        imageUrl,
		DisplayUrl: imageUrl,
		Mime:       "image/" + fileExtension,
	}
	slog.LogAttrs(context.Background(), slog.LevelInfo, "Image uploaded", slog.String("url", imageUrl))
	return c.JSON(&response)

}
