package tests

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/creasty/defaults"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/hetiansu5/urlquery"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

// Tester is a struct that mocks a request user
type Tester struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

var App *fiber.App

var (
	DefaultTester = &Tester{}
	UserTester    = &Tester{ID: 1}
	AdminTester   = &Tester{ID: 2}
	OtherTester   = map[int]*Tester{
		0: DefaultTester,
		1: UserTester,
		2: AdminTester,
	} // map[userID]Tester
)

// RegisterApp registers the fiber app to the package
// It should be called before any test
func RegisterApp(app *fiber.App) {
	App = app
}

// RequestConfig is a struct that contains the config of a request
type RequestConfig struct {
	RequestHeaders map[string]string
	RequestQuery   any
	RequestBody    any
	ResponseModel  any
	ExpectedBody   string
	ContentType    string `default:"application/json"`
}

func (tester *Tester) Request(t assert.TestingT, method string, route string, status int, config RequestConfig) {
	var requestData []byte
	var err error

	// set default values to config
	err = defaults.Set(&config)
	assert.Nilf(t, err, "set default values to config")

	model := config.ResponseModel

	// construct request
	if config.RequestQuery != nil {
		queryData, err := urlquery.Marshal(config.RequestQuery)
		assert.Nilf(t, err, "encode request query")
		route += "?" + string(queryData)
	}
	if config.RequestBody != nil {
		if data, ok := config.RequestBody.(string); ok {
			requestData = []byte(data)
		} else if config.ContentType == "application/json" {
			requestData, err = json.Marshal(config.RequestBody)
		} else if config.ContentType == "application/x-www-form-urlencoded" {
			requestData, err = urlquery.Marshal(config.RequestBody)
		} else {
			err = errors.New("unsupported content type: " + config.ContentType)
		}

		assert.Nilf(t, err, "encode request body")
	}
	req, err := http.NewRequest(
		method,
		route,
		bytes.NewBuffer(requestData),
	)
	assert.Nilf(t, err, "constructs http request")
	req.Header.Add("Content-Type", config.ContentType)
	if tester.Token != "" {
		req.Header.Add("Authorization", "Bearer "+tester.Token)
	}
	if config.RequestHeaders != nil {
		for key, value := range config.RequestHeaders {
			req.Header.Add(key, value)
		}
	}

	res, err := App.Test(req, -1)
	assert.Nilf(t, err, "perform request")
	assert.Equalf(t, status, res.StatusCode, "status code")

	responseBody, err := io.ReadAll(res.Body)
	assert.Nilf(t, err, "decode response")

	if res.StatusCode >= 400 {
		log.Print(string(responseBody))
	}

	if config.ExpectedBody != "" {
		assert.Equalf(t, config.ExpectedBody, string(responseBody), "response body")
	}
	if model != nil {
		err = json.Unmarshal(responseBody, model)
		assert.Nilf(t, err, "decode response")
	}
}

func (tester *Tester) Get(t assert.TestingT, route string, status int, config RequestConfig) {
	tester.Request(t, "GET", route, status, config)
}

func (tester *Tester) Post(t assert.TestingT, route string, status int, config RequestConfig) {
	tester.Request(t, "POST", route, status, config)
}

func (tester *Tester) Put(t assert.TestingT, route string, status int, config RequestConfig) {
	tester.Request(t, "PUT", route, status, config)
}

func (tester *Tester) Patch(t assert.TestingT, route string, status int, config RequestConfig) {
	tester.Request(t, "PATCH", route, status, config)
}

func (tester *Tester) Delete(t assert.TestingT, route string, status int, config RequestConfig) {
	tester.Request(t, "DELETE", route, status, config)
}
