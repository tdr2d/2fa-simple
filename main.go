package main

import (
	"2fa-simple/handlers"
	"2fa-simple/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
}

func main() {
	conf := new(utils.Config)
	if err := cleanenv.ReadConfig("config.yml", conf); err != nil {
		panic(err)
	}
	conf.TemplateDir = "templates"
	store := session.New(session.Config{CookieHTTPOnly: true, CookieSameSite: "true", Expiration: time.Hour})
	handler := handlers.Handler{Conf: conf, Store: store}
	engine := html.New("./templates", ".html")
	engine.Reload(true)

	// Middlewares
	app := fiber.New(fiber.Config{Views: engine})
	// app.Use(func(c *fiber.Ctx) error {
	// 	c.Set("X-XSS-Protection", "1; mode=block")
	// 	c.Set("X-Content-Type-Options", "nosniff")
	// 	c.Set("X-Download-Options", "noopen")
	// 	c.Set("Strict-Transport-Security", "max-age=5184000")
	// 	c.Set("X-Frame-Options", "SAMEORIGIN")
	// 	c.Set("X-DNS-Prefetch-Control", "off")
	// 	return c.Next()
	// })
	app.Use(recover.New())
	app.Use(logger.New())

	// Routes
	app.Get("/login", handler.LoginGetHandler)
	app.Post("/login", handler.LoginPostHandler)
	app.Post("/login/resend", handler.LoginResendHandler)
	app.Get("/login-check/:code", handler.LoginCheckHandler)
	app.Get("/logout", handler.LogoutGetHandler)
	app.Post("/password-reset", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/*", handler.SpaGetHandler)
	app.Static("/2fa-web", "./2fa-web")
	app.Listen(":3000")
}
