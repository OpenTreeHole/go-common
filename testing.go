package common

import (
	"bytes"
	"encoding/json"
	"github.com/creasty/defaults"
	"github.com/gofiber/fiber/v2"
	"github.com/hetiansu5/urlquery"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
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

// RegisterApp registers the fiber app to the common package
// It should be called before any test
func RegisterApp(app *fiber.App) {
	App = app
}

// RequestConfig is a struct that contains the config of a request
type RequestConfig struct {
	Method         string            `default:"GET"`
	Route          string            `default:"/"`
	ExpectedStatus int               `default:"200"`
	RequestHeaders map[string]string `default:"-"`
	RequestQuery   any               `default:"-"`
	RequestBody    any               `default:"-"`
	ResponseModel  any               `default:"-"`
	ExpectedBody   string            `default:"-"`
}

func (tester *Tester) Request(t assert.TestingT, config RequestConfig) {
	var requestData []byte
	var err error

	// set default values to config
	err = defaults.Set(&config)
	assert.Nilf(t, err, "set default values to config")

	method := config.Method
	route := config.Route
	statusCode := config.ExpectedStatus
	model := config.ResponseModel

	// construct request
	if config.RequestQuery != nil {
		queryData, err := urlquery.Marshal(config.RequestQuery)
		assert.Nilf(t, err, "encode request query")
		route += "?" + string(queryData)
	}
	if config.RequestBody != nil {
		requestData, err = json.Marshal(config.RequestBody)
		assert.Nilf(t, err, "encode request body")
	}
	req, err := http.NewRequest(
		method,
		route,
		bytes.NewBuffer(requestData),
	)
	assert.Nilf(t, err, "constructs http request")
	req.Header.Add("Content-Type", "application/json")
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
	assert.Equalf(t, statusCode, res.StatusCode, "status code")

	responseBody, err := io.ReadAll(res.Body)
	assert.Nilf(t, err, "decode response")

	if res.StatusCode >= 400 {
		log.Print(string(responseBody))
	} else {
		if config.ExpectedBody != "" {
			assert.Equalf(t, config.ExpectedBody, string(responseBody), "response body")
		}
		if model != nil {
			err = json.Unmarshal(responseBody, model)
			assert.Nilf(t, err, "decode response")
		}
	}
}

func (tester *Tester) Get(t assert.TestingT, config RequestConfig) {
	config.Method = "GET"
	tester.Request(t, config)
}

func (tester *Tester) Post(t assert.TestingT, config RequestConfig) {
	config.Method = "POST"
	if config.ExpectedStatus == 0 {
		config.ExpectedStatus = fiber.StatusCreated
	}
	tester.Request(t, config)
}

func (tester *Tester) Put(t assert.TestingT, config RequestConfig) {
	config.Method = "PUT"
	tester.Request(t, config)
}

func (tester *Tester) Patch(t assert.TestingT, config RequestConfig) {
	config.Method = "PATCH"
	tester.Request(t, config)
}

func (tester *Tester) Delete(t assert.TestingT, config RequestConfig) {
	config.Method = "DELETE"
	tester.Request(t, config)
}
