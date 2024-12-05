package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	. "github.com/opentreehole/backend/image_hosting/model"
	. "github.com/opentreehole/backend/image_hosting/schema"
	"io"
	"log"
	"path/filepath"
	"time"
)

func GenerateUniqid() string {
	// Get the current timestamp in microseconds
	now := time.Now().UnixMicro()

	// Generate a random 6-byte (12-character) string
	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err) // Handle error if random generation fails
	}
	randomSuffix := hex.EncodeToString(randomBytes) // Convert to hex string

	// Combine the timestamp and random suffix
	return fmt.Sprintf("%x%s", now, randomSuffix)
}

// to-do-list
// uploadImage里面生成一下basename的逻辑（根据当前时间生成）
// 	研究具体逻辑，见notion旦夕file，查找basename / filename：https://github.com/lsky-org/lsky-pro/blob/master/app/Services/ImageService.php#L185
// https://github.com/lsky-org/lsky-pro/blob/master/config/convention.php#L113

// 本函数是一个上传接口，需要实现把图片传到数据库，由此需要为图片生成一个basename，basename是lsky pro里面的命名规范（开源，直接从github代码里看）
// CheveretoUploadResponse是返回给前端的一个Response，见api.go里会把这个return回前端，前端会用里面的height/url等作展示，因此仍需要继续使用这个response，在函数中往这个response存信息即可（以下是上传成功的话要传的信息）：

func UploadImage(c *fiber.Ctx) error {
	log.Println("uploading image")
	var response CheveretoUploadResponse

	// 传jpg没问题，但传png有问题（有可能和下文代码有关）
	file, err := c.FormFile("source")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Bad Request: Unable to retrieve file.",
		})
	}
	fileExtension := filepath.Ext(file.Filename)
	fileContent, err := file.Open()
	identifier := GenerateUniqid()
	if err != nil {
		return err
	}

	// 需要存进数据库的变量
	content, err := io.ReadAll(fileContent)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Something goes wrong in file data",
		})
	}
	directUrl := time.Now().Format("2006/01/02/") + identifier + fileExtension
	imageName := identifier
	uploadedImage := &ImageTable{
		BaseName:      imageName,
		SavingTime:    time.Now(),
		ImageType:     fileExtension,
		ImageFileData: content,
	}
	err = DB.Create(&uploadedImage).Error
	println("upload image")
	if err != nil {
		panic(err)
		return err
	}

	response.StatusCode = 200
	response.StatusTxt = "Upload Success"
	response.Image = CheveretoImageInfo{
		// strings.TrimSuffix(lskyUploadResponse.Data.Name, filepath.Ext(lskyUploadResponse.Data.Name)),
		Name:      imageName,
		Extension: fileExtension,
		// 	Md5:        lskyUploadResponse.Data.Md5,
		// lskyUploadResponse.Data.Name
		Filename: imageName + fileExtension,
		// 	Mime:       lskyUploadResponse.Data.Mimetype,
		Url:        directUrl,
		DisplayUrl: directUrl,
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
