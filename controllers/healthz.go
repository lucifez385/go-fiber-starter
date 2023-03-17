package controllers

import (
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
)

type healthzController interface {
	GetHealthz(c *fiber.Ctx) error
}

type HealthzController struct {
	Log utils.ILogging
}

func NewHealthzController() healthzController {
	return &HealthzController{
		Log: &utils.Logging{},
	}
}

func (ctr *HealthzController) GetHealthz(c *fiber.Ctx) error {
	code := fiber.StatusOK
	ctr.Log.Info(c, code, "OK", nil)
	return c.Status(code).SendString("OK")
}
