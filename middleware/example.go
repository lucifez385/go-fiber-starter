package middleware

import (
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
)

type exampleMiddleware interface {
	GetAuth(c *fiber.Ctx) error
}

type ExampleMiddleware struct {
	Log utils.ILogging
}

func NewExampleMiddleware() exampleMiddleware {
	return &ExampleMiddleware{
		Log: &utils.Logging{},
	}
}

func (ctr *ExampleMiddleware) GetAuth(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	if authorization == "" {
		err := fiber.ErrUnauthorized
		ctr.Log.Error(c, err.Code, err.Error(), err)
		return c.Status(err.Code).JSON(err)
	}
	c.Locals("name", "user.name")
	c.Locals("authorization", authorization)
	return c.Next()
}
