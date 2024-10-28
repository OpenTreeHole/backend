package service

import (
	"fmt"
	"github.com/opentreehole/backend/image_hosting/schema"
	"io"
	"mime/multipart"
)

func UploadImage(file *multipart.FileHeader, response *schema.CheveretoUploadResponse) error {

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

	fmt.Println(content)
	return nil
}

// // --------------------------- 这里开始就是lsky的图床 -------------------------

// fileData貌似是用于给lsky设置的
// fileData := &fiber.FormFile{
// 	Fieldname: "file",
// 	Name:      file.Filename,
// 	Content:   content,
// }

// agent := LskyBaseAgent(fiber.AcquireAgent(), fiber.MethodPost, "/upload").
// 	Set("Authorization", "Bearer "+token).
// 	ContentType(fiber.MIMEMultipartForm).
// 	FileData(fileData).MultipartForm(nil)
// defer fiber.ReleaseAgent(agent)

// for {
// 	agent.Set("Authorization", "Bearer "+token)
// 	if err = agent.Parse(); err != nil {
// 		return err
// 	}

// 	code, body, errs := agent.Bytes()
// 	if len(errs) != 0 {
// 		return errs[0]
// 	}
// 	if code == 200 {
// 		var lskyUploadResponse LskyUploadResponse
// 		err = json.Unmarshal(body, &lskyUploadResponse)
// 		if err != nil {
// 			return err
// 		}

// 		// url transform to direct url
// 		urlRaw := lskyUploadResponse.Data.Links.Url
// 		urlData, err := url.ParseRequestURI(urlRaw)
// 		if err != nil {
// 			return err
// 		}
// 		if Config.HostRewrite != "" {
// 			urlData.Host = Config.HostRewrite
// 		} else {
// 			urlData.Host = ProxyUrlData.Host
// 		}
// 		urlData.Scheme = ProxyUrlData.Scheme
// 		directUrl := urlData.String()
// 		log.Printf("image upload: %v\n", directUrl)

// 		return nil
// 	} else if code == 401 {
// 		// refresh token
// 		newToken := GetToken() // maybe another coroutine refresh the token
// 		if token != newToken {
// 			// another coroutine refresh the token
// 			token = newToken
// 		} else {
// 			// this coroutine refresh the token
// 			token, err = LskyRefreshToken()
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	} else {
// 		message := fmt.Sprintf(`{"code": %v}`, code)
// 		return fiber.NewError(fiber.StatusInternalServerError, message)
// 	}
// }
