package api

import (
	"github.com/gofiber/fiber/v2"
	. "github.com/opentreehole/backend/image_hosting/model"
	. "github.com/opentreehole/backend/image_hosting/schema"
	. "github.com/opentreehole/backend/image_hosting/utils"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// 把各种error都处理一下（including other file）
// 运行时的log
// 缩图，url放进thumb的url字段
// main.go里的中间件
// 写注释
// 修改不清楚的变量名

func UploadImage(c *fiber.Ctx) error {
	// response to frontend
	var response CheveretoUploadResponse
	// the file uploaded by user
	file, err := c.FormFile("source")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Bad Request: Unable to retrieve file or no file has been sent.",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Bad Request: Unable to open file.",
		})
	}

	imageData, err := io.ReadAll(fileContent)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "The uploaded file has some problems",
		})
	}

	fileExtension := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
	imageIdentifier := GenerateIdentifier()

	imageUrl := Hostname + "/api/i/" + time.Now().Format("2006/01/02/") + imageIdentifier + "." + fileExtension
	uploadedImage := &ImageTable{
		BaseName:      imageIdentifier,
		ImageType:     fileExtension,
		ImageFileData: imageData,
	}
	err = DB.Create(&uploadedImage).Error

	if err != nil {
		return err
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
	}
	log.Printf("image upload: %v\n", imageUrl)
	return c.JSON(&response)

}

func GetImage(c *fiber.Ctx) error {
	// to access the image in database
	// year := c.Params("year")
	// month := c.Params("month")
	// day := c.Params("day")
	// imageType := strings.Split(c.Params("imageIdentifier"), ".")[1]

	var image ImageTable
	imageIdentifier := strings.Split(c.Params("identifier"), ".")[0]
	DB.First(&image, "base_name = ?", imageIdentifier)
	// browser will automatically transform the BLOB data to an image (no matter what extension)
	return c.Send(image.ImageFileData)
}
