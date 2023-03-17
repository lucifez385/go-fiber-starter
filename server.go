package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-fiber-starter/routers"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	logging := &utils.Logging{}

	handleArgs()

	app := fiber.New(fiber.Config{
		BodyLimit: 400 * 1024 * 1024,
	})
	app.Use(cors.New())
	app.Use(recover.New())

	routers.Route(app)

	go func() {
		if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
			fmt.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	signals := <-c
	switch signals {
	case os.Interrupt:
		logging.InfoWithoutCtx(200, "Got SIGINT...", "Got SIGINT...")
	case syscall.SIGTERM:
		logging.InfoWithoutCtx(200, "Got SIGTERM...", "Got SIGTERM...")
	}

	if err := app.Shutdown(); err != nil {
		logging.ErrorWithoutCtx(500, err.Error(), err)
	} else {
		logging.InfoWithoutCtx(200, "gracefully shutting down", "gracefully shutting down")
	}
}

func handleArgs() {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "healthz":
			utils.HealthCheck()
		}
		os.Exit(0)
	}
}
