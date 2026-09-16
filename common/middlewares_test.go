package common

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareCustomLoggerDoesNotLogRequestBody(t *testing.T) {
	var logOutput bytes.Buffer
	previousLogger := Logger
	Logger = slog.New(slog.NewJSONHandler(&logOutput, nil))
	t.Cleanup(func() { Logger = previousLogger })

	app := fiber.New()
	app.Use(MiddlewareCustomLogger)
	app.Post("/upload", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	const sensitiveBody = "private-image-binary-data"
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("source", "private.heic")
	require.NoError(t, err)
	_, err = part.Write([]byte(sensitiveBody))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	req := httptest.NewRequest("POST", "/upload", requestBody)
	req.Header.Set(fiber.HeaderContentType, writer.FormDataContentType())

	response, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(logOutput.Bytes(), &entry))
	require.NotContains(t, entry, "body")
	require.NotContains(t, logOutput.String(), sensitiveBody)
	require.Equal(t, "POST", entry["method"])
	require.Equal(t, "/upload", entry["origin_url"])
}
