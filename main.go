package main

import (
	"2fa-simple/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
}

func main() {
	engine := html.New("./templates", ".html")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(recover.New())
	app.Use(logger.New())
	// store := session.New()

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{"Title": "Login"}, "layout")
	})
	app.Post("/login", handler.LoginHandler)

	app.Post("/forgot-password", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Static("/", "./public")
	app.Listen(":3000")
}
