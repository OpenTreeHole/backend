package api

import (
	"github.com/gofiber/fiber/v2"
	. "github.com/opentreehole/backend/image_hosting/model"
	. "github.com/opentreehole/backend/image_hosting/schema"
	"io"
	"log"
	"time"
)

func UploadImage(c *fiber.Ctx) error {
	log.Println("uploading image")
	var response CheveretoUploadResponse

	// 传jpg没问题，但传png有问题（有可能和下文代码有关）
	file, err := c.FormFile("source")
	if err != nil {
		log.Println(err)
		return err
	}

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
		log.Println(err)
		return err
	}

	// 需要存进数据库的变量
	content, err := io.ReadAll(fileContent)

	if err != nil {
		log.Println(err)
		return err
	}
	println(content)
	uploadedImage := &ImageTable{
		BaseName:      "example_image",
		SavingTime:    time.Now(),
		ImageType:     "jpeg",
		ImageFileData: content,
	}
	err = DB.Create(&uploadedImage).Error
	println("upload image")
	if err != nil {
		panic(err)
		return err
	}

	//
	// 	courseGroup = request.ToCourseGroupModel()
	// 	err = DB.Create(&courseGroup).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	//
	// course = request.ToModel(courseGroup.ID)
	// course.CourseGroup = courseGroup
	// err = course.Create()
	// if err != nil {
	// return err
	// }
	//
	// return c.JSON(new(CourseV1Response).FromModel(user, course))

	return c.JSON(&response)

}

func GetImage(c *fiber.Ctx) error {
	// to access the image in database
	// year := c.Params("year")
	// month := c.Params("month")
	// day := c.Params("day")
	//
	// identifierName := strings.Split(c.Params("identifier"), ".")[0]
	// imageType := strings.Split(c.Params("identifier"), ".")[1]

	// 有了以上，可以用它们组合起来来在数据库中查找图片数据
	var image ImageTable
	DB.First(&image, 1)
	return c.Send(image.ImageFileData)
}
