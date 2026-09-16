package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/backend/common"
	"github.com/opentreehole/backend/image_hosting/config"
	"github.com/opentreehole/backend/image_hosting/model"
	"github.com/opentreehole/backend/image_hosting/schema"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUploadImageAcceptsHEICAndJPEG(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), common.GormConfig)
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&model.ImageTable{}))

	previousDB := model.DB
	model.DB = database
	t.Cleanup(func() { model.DB = previousDB })

	previousHostName := viper.GetString(config.EnvHostName)
	viper.Set(config.EnvHostName, "https://images.example.com")
	t.Cleanup(func() { viper.Set(config.EnvHostName, previousHostName) })

	app := fiber.New(fiber.Config{ErrorHandler: common.ErrorHandler})
	app.Post("/api/uploadImage", UploadImage)

	testCases := []struct {
		name        string
		filename    string
		contentType string
		data        []byte
		extension   string
	}{
		{
			name:        "Android HEIC",
			filename:    "android-upload.heic",
			contentType: "image/heic",
			data:        []byte("\x00\x00\x00\x18ftypheic-test-image"),
			extension:   "heic",
		},
		{
			name:        "JPEG",
			filename:    "regular-upload.jpeg",
			contentType: "image/jpeg",
			data:        []byte("\xff\xd8\xff\xe0-test-image"),
			extension:   "jpeg",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			requestBody, contentType := multipartImageBody(
				t,
				testCase.filename,
				testCase.contentType,
				testCase.data,
			)
			req := httptest.NewRequest("POST", "/api/uploadImage", requestBody)
			req.Header.Set(fiber.HeaderContentType, contentType)

			response, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusOK, response.StatusCode)

			var uploadResponse schema.CheveretoUploadResponse
			require.NoError(t, json.NewDecoder(response.Body).Decode(&uploadResponse))
			require.Equal(t, 200, uploadResponse.StatusCode)
			require.Equal(t, testCase.extension, uploadResponse.Image.Extension)
			require.Equal(t, testCase.contentType, uploadResponse.Image.Mime)

			var storedImage model.ImageTable
			require.NoError(t, database.First(&storedImage, "original_file_name = ?", testCase.filename).Error)
			require.Equal(t, testCase.extension, storedImage.ImageType)
			require.Equal(t, testCase.data, storedImage.ImageFileData)
		})
	}
}

func multipartImageBody(t *testing.T, filename, contentType string, data []byte) (*bytes.Buffer, string) {
	t.Helper()

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="source"; filename="`+filename+`"`)
	header.Set("Content-Type", contentType)
	part, err := writer.CreatePart(header)
	require.NoError(t, err)
	_, err = part.Write(data)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	return requestBody, writer.FormDataContentType()
}
