package controllers

import (
	"go-fiber-starter/types"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
)

type exampleController interface {
	GetWithMiddleware(c *fiber.Ctx) error
	PostWithValidator(c *fiber.Ctx) error
	// GetCansExample(c *fiber.Ctx) error
	// Logout(c *fiber.Ctx) error
}

type ExampleController struct {
	Log utils.ILogging
}

type ExampleData struct {
	name          string
	authorization string
}

func NewExampleController() exampleController {
	return &ExampleController{
		Log: &utils.Logging{},
	}
}

func (ctr *ExampleController) GetWithMiddleware(c *fiber.Ctx) error {
	name := c.Locals("name").(string)
	authorization := c.Locals("authorization").(string)
	status := fiber.StatusOK
	data := ExampleData{
		name:          name,
		authorization: authorization,
	}
	ctr.Log.Info(c, status, "Operation completed", data)
	return c.Status(status).JSON(data)
}

func (ctr *ExampleController) PostWithValidator(c *fiber.Ctx) error {
	authorization := c.Locals("authorization").(string)
	status := fiber.StatusOK

	body := new(types.ExampleBody)
	if err := c.BodyParser(body); err != nil {
		status := fiber.StatusInternalServerError
		ctr.Log.Error(c, status, err.Error(), err)
		return c.Status(status).JSON(err)
	}

	data := ExampleData{
		name:          body.Name,
		authorization: authorization,
	}
	ctr.Log.Info(c, status, "Operation completed", data)
	return c.Status(status).JSON(data)
}
