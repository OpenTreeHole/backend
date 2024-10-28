package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/backend/image_hosting/schema"
	"io"
)

func UploadImage(c *fiber.Ctx) error {
	var response schema.CheveretoUploadResponse
	file, err := c.FormFile("source")
	if err != nil {
		return err
	}

	//err = service.UploadImage(file, &response)

	// 本函数是一个上传接口，需要实现把图片传到数据库，由此需要为图片生成一个basename，basename是lsky pro里面的命名规范（开源，直接从github代码里看）
	// CheveretoUploadResponse是返回给前端的一个Response，见api.go里会把这个return回前端，前端会用里面的height/url等作展示，因此仍需要继续使用这个response，在函数中往这个response存信息即可（以下是上传成功的话要传的信息）：
	// 		response.StatusCode = 200
	// 		response.StatusTxt = "Upload Success"
	// 		response.Image = CheveretoImageInfo{
	// 			Name:       strings.TrimSuffix(lskyUploadResponse.Data.Name, filepath.Ext(lskyUploadResponse.Data.Name)),
	// 			Extension:  lskyUploadResponse.Data.Extension,
	// 			Md5:        lskyUploadResponse.Data.Md5,
	// 			Filename:   lskyUploadResponse.Data.Name,
	// 			Mime:       lskyUploadResponse.Data.Mimetype,
	// 			Url:        directUrl,
	// 			DisplayUrl: directUrl,
	// 		}

	fileContent, err := file.Open()
	if err != nil {
		return err
	}

	content, err := io.ReadAll(fileContent)

	if err != nil {
		return err
	}
	println(content)

	return c.JSON(&response)

}

func GetImage(c *fiber.Ctx) error {
	// to access the image in database
	//year := c.Params("year")
	//month := c.Params("month")
	//day := c.Params("day")
	//identifier := c.Params("identifier")
	//
	//identifier = strings.Split(identifier, ".")[0]
	//imageType = strings.Split(identifier, ".")[1]

	// 有了以上，可以用它们组合起来来在数据库中查找图片数据
	// 函数最后return image的BLOB信息，前端会自动解析为一张图片的
	return nil
}
