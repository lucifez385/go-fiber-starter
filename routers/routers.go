package routers

import (
	"go-fiber-starter/controllers"

	"go-fiber-starter/middleware"
	"go-fiber-starter/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {

	// Middleware configuration
	exampleMiddleware := middleware.NewExampleMiddleware()

	// Controller configuration
	healthzController := controllers.NewHealthzController()
	exampleController := controllers.NewExampleController()

	// Router definition
	app.Get("/healthz", healthzController.GetHealthz)

	exampleV1 := app.Group("example/v1")
	exampleV1.Get("/middleware", exampleMiddleware.GetAuth, exampleController.GetWithMiddleware)
	exampleV1.Post("/validator", validator.Example, exampleMiddleware.GetAuth, exampleController.PostWithValidator)

}
