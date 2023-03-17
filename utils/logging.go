package utils

import (
	"encoding/json"
	"fmt"
	"go-fiber-starter/constants"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var FmtPrintLn = fmt.Println

type LogAttributes struct {
	HttpMethod            string `json:"http.method"`
	HttpUrl               string `json:"http.url"`
	HttpPath              string `json:"http.path"`
	HttpHost              string `json:"http.host"`
	HttpScheme            string `json:"http.scheme"`
	HttpStatusCode        int    `json:"http.status_code"`
	HttpFlavor            string `json:"http.flavor"`
	UserAgent             string `json:"http.user_agent"`
	ContentLength         int    `json:"http.request_content_length"`
	ResponseContentLength int    `json:"http.response_content_length"`
}
type LogFormat struct {
	Timestamp    string `json:"timestamp"`
	SeverityText string `json:"severityText"`
	Body         string `json:"body"`
	// Resource     LogResource   `json:"resource"`
	Attributes LogAttributes `json:"attributes"`
	Kind       string        `json:"kind"`
}

type ILogging interface {
	Emit(ctx *fiber.Ctx, status int, message string, rawData interface{}, currentLogLevel string)
	Debug(ctx *fiber.Ctx, status int, message string, rawData interface{})
	Info(ctx *fiber.Ctx, status int, message string, rawData interface{})
	Warning(ctx *fiber.Ctx, status int, message string, rawData interface{})
	Error(ctx *fiber.Ctx, status int, message string, rawData interface{})
	Fatal(ctx *fiber.Ctx, status int, message string, rawData interface{})

	EmitWithoutCtx(status int, message string, rawData interface{}, currentLogLevel string)
	DebugWithoutCtx(status int, message string, rawData interface{})
	InfoWithoutCtx(status int, message string, rawData interface{})
	WarningWithoutCtx(status int, message string, rawData interface{})
	ErrorWithoutCtx(status int, message string, rawData interface{})
	FatalWithoutCtx(status int, message string, rawData interface{})
}

type Logging struct{}

func (c *Logging) Emit(ctx *fiber.Ctx, status int, message string, data interface{}, currentLogLevel string) {
	logLevel := "info"
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	logEmit := false
	for _, key := range constants.LOG_LEVEL_FILTER[logLevel] {
		if key == currentLogLevel {
			logEmit = true
			break
		}
	}
	if logEmit {
		logBytes, err := json.Marshal(data)
		if err != nil {
			panic("Logging error")
		}
		flavorArr := strings.Split(string(ctx.Request().Header.Protocol()), "/")
		flavor := ""
		if len(flavorArr) >= 2 {
			flavor = strings.ToUpper(flavorArr[1])
		}
		reqContentLength := ctx.Request().Header.ContentLength()
		if reqContentLength < 0 {
			reqContentLength = 0
		}

		attributes := LogAttributes{
			HttpMethod:            string(ctx.Request().Header.Method()),
			HttpUrl:               string(ctx.BaseURL()) + string(ctx.Request().RequestURI()),
			HttpPath:              string(ctx.Request().RequestURI()),
			HttpHost:              string(ctx.Request().Host()),
			HttpScheme:            string(ctx.Request().URI().Scheme()),
			HttpStatusCode:        status,
			HttpFlavor:            flavor,
			UserAgent:             string(ctx.Request().Header.UserAgent()),
			ContentLength:         reqContentLength,
			ResponseContentLength: len(logBytes),
		}

		timeNow := time.Now()
		timeInSec := fmt.Sprint(timeNow.Unix())
		timeInMilli := fmt.Sprint(timeNow.UnixNano() / 1000000)
		millSecUnit := timeInMilli[10:13]
		formattedTime := timeInSec + "." + millSecUnit

		logInfo := LogFormat{
			Timestamp:    formattedTime,
			SeverityText: strings.ToUpper(currentLogLevel),
			Body:         message,
			// Resource: LogResource{
			// 	ServiceName:    os.Getenv("APP_NAME"),
			// 	ServiceVersion: os.Getenv("VERSION"),
			// },
			Attributes: attributes,
			Kind:       "application",
		}
		js, _ := json.Marshal(logInfo)

		if _, err := FmtPrintLn(string(js)); err != nil {
			panic(err.Error())
		}
	}
}

func (c *Logging) EmitWithoutCtx(status int, message string, data interface{}, currentLogLevel string) {
	timeNow := time.Now()
	timeInSec := fmt.Sprint(timeNow.Unix())
	timeInMilli := fmt.Sprint(timeNow.UnixNano() / 1000000)
	millSecUnit := timeInMilli[10:13]
	formattedTime := timeInSec + "." + millSecUnit

	logByte, err := json.Marshal(data)
	if err != nil {
		panic("Logging error")
	}

	attributes := LogAttributes{
		HttpStatusCode:        status,
		ResponseContentLength: len(logByte),
	}

	logInfo := LogFormat{
		Timestamp:    formattedTime,
		SeverityText: strings.ToUpper(currentLogLevel),
		Body:         message,
		Attributes:   attributes,
		Kind:         "application",
	}
	js, _ := json.Marshal(logInfo)

	if _, err := FmtPrintLn(string(js)); err != nil {
		panic(err.Error())
	}
}

func (c *Logging) Debug(ctx *fiber.Ctx, status int, message string, rawData interface{}) {
	c.Emit(ctx, status, message, rawData, constants.LOG_LEVEL.DEBUG)
}

func (c *Logging) Info(ctx *fiber.Ctx, status int, message string, rawData interface{}) {
	c.Emit(ctx, status, message, rawData, constants.LOG_LEVEL.INFO)
}

func (c *Logging) Warning(ctx *fiber.Ctx, status int, message string, rawData interface{}) {
	c.Emit(ctx, status, message, rawData, constants.LOG_LEVEL.WARNING)
}

func (c *Logging) Error(ctx *fiber.Ctx, status int, message string, rawData interface{}) {
	c.Emit(ctx, status, message, rawData, constants.LOG_LEVEL.ERROR)
}

func (c *Logging) Fatal(ctx *fiber.Ctx, status int, message string, rawData interface{}) {
	c.Emit(ctx, status, message, rawData, constants.LOG_LEVEL.FATAL)
}

func (c *Logging) DebugWithoutCtx(status int, message string, rawData interface{}) {
	c.EmitWithoutCtx(status, message, rawData, constants.LOG_LEVEL.DEBUG)
}

func (c *Logging) InfoWithoutCtx(status int, message string, rawData interface{}) {
	c.EmitWithoutCtx(status, message, rawData, constants.LOG_LEVEL.INFO)
}

func (c *Logging) WarningWithoutCtx(status int, message string, rawData interface{}) {
	c.EmitWithoutCtx(status, message, rawData, constants.LOG_LEVEL.WARNING)
}

func (c *Logging) ErrorWithoutCtx(status int, message string, rawData interface{}) {
	c.EmitWithoutCtx(status, message, rawData, constants.LOG_LEVEL.ERROR)
}

func (c *Logging) FatalWithoutCtx(status int, message string, rawData interface{}) {
	c.EmitWithoutCtx(status, message, rawData, constants.LOG_LEVEL.FATAL)
}
