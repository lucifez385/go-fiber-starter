package validator

import (
	"go-fiber-starter/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ExamplePayload struct {
	Name string `validate:"required,max=100"`
}

func Example(c *fiber.Ctx) error {
	logging := &utils.Logging{}
	body := new(ExamplePayload)
	if err := c.BodyParser(body); err != nil {
		code := fiber.StatusInternalServerError
		logging.Error(c, code, err.Error(), err)
		return c.Status(code).JSON(err)
	}
	validate := validator.New()
	err := validate.Struct(body)

	errs := validationErrorFormat(err)
	if len(errs) > 0 {
		err := fiber.ErrUnprocessableEntity
		logging.Error(c, err.Code, err.Error(), err)
		return c.Status(err.Code).JSON(errs)
	}
	return c.Next()
}
