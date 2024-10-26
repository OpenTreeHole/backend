package main

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func UploadImage(c *fiber.Ctx) error {
	var response CheveretoUploadResponse
	file, err := c.FormFile("source")
	if err != nil {
		return err
	}
	err = ProxyUploadImage(file, &response)
	if err != nil {
		return err
	}
	return c.JSON(&response)
}

func GetImage(c *fiber.Ctx) error {
	// to access the image in database
	year := c.Params("year")
	month := c.Params("month")
	day := c.Params("day")
	identifier := c.Params("identifier")

	identifier = strings.Split(identifier, ".")[0]
	imageType = strings.Split(identifier, ".")[1]

	// 有了以上，可以用它们组合起来来在数据库中查找图片数据
	// 函数最后return image的BLOB信息，前端会自动解析为一张图片的
	return nil
}
