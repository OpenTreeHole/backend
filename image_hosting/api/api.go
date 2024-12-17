package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/backend/common"
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
	// the file uploaded by user  in the request body with the form-data key "source"
	file, err := c.FormFile("source")
	if err != nil {
		log.Println(err)
		return common.BadRequest("No file uploaded")
	}

	fileSize := file.Size
	maxSize := 10 * 1024 * 1024 // file should <= 10MB
	if int(fileSize) > maxSize {
		return common.BadRequest("File size is too large")
	}

	fileExtension := strings.TrimPrefix(filepath.Ext(file.Filename), ".")

	if !IsAllowedExtension(fileExtension) {
		log.Println("Error: File type not allowed.")
		return common.BadRequest("File type not allowed.")
	}

	fileContent, err := file.Open()
	if err != nil {
		log.Println(err)
		return common.BadRequest("The uploaded file has some problems")
	}

	imageData, err := io.ReadAll(fileContent)
	if err != nil {
		log.Println(err)
		return common.BadRequest("The uploaded file has some problems")
	}

	imageIdentifier, err := GenerateIdentifier()
	if err != nil {
		log.Println(err)
		return common.InternalServerError("Cannot generate image identifier")
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
		return common.BadRequest("Cannot find the image")
	}
	log.Println("get image successfully")
	// browser will automatically transform the BLOB data to an image (no matter what extension)
	return c.Send(image.ImageFileData)
}
