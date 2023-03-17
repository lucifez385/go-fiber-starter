package utils_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go-fiber-starter/constants"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/go-playground/assert.v1"

	"github.com/stretchr/testify/suite"
	"github.com/tkuchiki/faketime"
)

type TestSuiteLogging struct {
	suite.Suite
	LogJson utils.LogFormat
}

func (suite *TestSuiteLogging) SetupTest() {
	utils.FmtPrintLn = fmt.Println
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("APP_NAME", "editor-api")
	os.Setenv("VERSION", "v0.0.0-beta.1")

	attributes := utils.LogAttributes{
		HttpMethod:            "GET",
		HttpUrl:               "http://example.com/test/logging",
		HttpPath:              "/test/logging",
		HttpHost:              "example.com",
		HttpScheme:            "http",
		HttpStatusCode:        200,
		HttpFlavor:            "1.1",
		UserAgent:             "",
		ContentLength:         0,
		ResponseContentLength: 4,
	}
	logInfo := utils.LogFormat{
		Timestamp:    "1257894000.000",
		SeverityText: "INFO",
		Body:         "Emit Logging Testing",
		// Resource: utils.LogResource{
		// 	ServiceName:    "editor-api",
		// 	ServiceVersion: "v0.0.0-beta.1",
		// },
		Attributes: attributes,
		Kind:       "application",
	}
	suite.LogJson = logInfo
}

func (suite *TestSuiteLogging) TestEmitLog() {
	app := fiber.New()
	t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(t)
	defer f.Undo()
	f.Do()

	var logResult interface{}
	utils.FmtPrintLn = func(a ...interface{}) (n int, err error) {
		logResult = a[0]
		return 0, nil
	}

	app.Get("/test/logging", func(c *fiber.Ctx) error {
		log := utils.Logging{}
		var rawData interface{}
		log.Emit(c, 200, "Emit Logging Testing", rawData, constants.LOG_LEVEL.INFO)
		return c.SendString("")
	})

	req := httptest.NewRequest("GET", "/test/logging", nil)

	if _, err := app.Test(req); err != nil {
		suite.Error(err)
	}
	lg := logResult.(string)
	result := utils.LogFormat{}
	if err := json.Unmarshal([]byte(lg), &result); err != nil {
		suite.Error(err)
	}
	suite.Equal(suite.LogJson, result)
}

func (suite *TestSuiteLogging) TestLogLevelUpThanCurrentLevel() {
	os.Setenv("LOG_LEVEL", "error")
	app := fiber.New()
	t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(t)
	defer f.Undo()
	f.Do()
	var logResult interface{}
	logResult = ""
	utils.FmtPrintLn = func(a ...interface{}) (n int, err error) {
		logResult = a[0]
		return 0, nil
	}
	app.Get("/test/logging", func(c *fiber.Ctx) error {
		log := utils.Logging{}
		var rawData interface{}
		log.Emit(c, 200, "Emit Logging Testing", rawData, constants.LOG_LEVEL.INFO)
		return c.SendString("")
	})

	req := httptest.NewRequest("GET", "/test/logging", nil)

	if _, err := app.Test(req); err != nil {
		suite.Error(err)
	}
	lg := logResult.(string)
	assert.Equal(suite.T(), "", lg)
}

func TestLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteLogging))
}
