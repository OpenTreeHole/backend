package api

import (
	"github.com/gofiber/fiber/v2"
	. "github.com/opentreehole/backend/image_hosting/config"
	. "github.com/opentreehole/backend/image_hosting/model"
	. "github.com/opentreehole/backend/image_hosting/schema"
	. "github.com/opentreehole/backend/image_hosting/utils"
	"io"
	"log"

	"path/filepath"
	"strings"
	"time"
)

func UploadImage(c *fiber.Ctx) error {
	log.Println("uploading image")
	// response to frontend
	var response CheveretoUploadResponse
	// the file uploaded by user
	file, err := c.FormFile("source")
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Bad Request: Unable to retrieve file or no file has been sent.",
		})
	}

	fileSize := file.Size
	maxSize := 10 * 1024 * 1024 // file should <= 10MB
	if int(fileSize) > maxSize {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"code":    413,
			"message": "File size exceeds the maximum limit of 10MB",
		})
	}

	fileExtension := strings.TrimPrefix(filepath.Ext(file.Filename), ".")

	if !IsAllowedExtension(fileExtension) {
		log.Println("Error: File type not allowed.")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Bad Request: File type not allowed.",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Bad Request: Unable to open file.",
		})
	}

	imageData, err := io.ReadAll(fileContent)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "The uploaded file has some problems",
		})
	}

	imageIdentifier, err := GenerateIdentifier()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Cannot generate image identifier",
		})
	}

	imageUrl := Config.HostName + "/api/i/" + time.Now().Format("2006/01/02/") + imageIdentifier + "." + fileExtension
	uploadedImage := &ImageTable{
		BaseName:      imageIdentifier,
		ImageType:     fileExtension,
		ImageFileData: imageData,
	}
	err = DB.Create(&uploadedImage).Error

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Database cannot store the image",
		})
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
	log.Printf("image is uploaded on: %v\n", imageUrl)
	return c.JSON(&response)

}

func GetImage(c *fiber.Ctx) error {
	// to access the image in database
	// year := c.Params("year")
	// month := c.Params("month")
	// day := c.Params("day")
	// imageType := strings.Split(c.Params("imageIdentifier"), ".")[1]
	log.Println("getting image")
	var image ImageTable
	imageIdentifier := strings.Split(c.Params("identifier"), ".")[0]
	err := DB.First(&image, "base_name = ?", imageIdentifier)
	if err.Error != nil {
		log.Println(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    404,
			"message": "Image not found",
		})
	}
	log.Println("get image successfully")
	// browser will automatically transform the BLOB data to an image (no matter what extension)
	return c.Send(image.ImageFileData)
}
