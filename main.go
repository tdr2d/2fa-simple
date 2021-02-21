package main

import (
	"2fa-simple/handlers"
	"2fa-simple/utils"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"github.com/gofiber/template/html"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
}

func main() {
	// Config
	conf := new(utils.Config)
	if err := cleanenv.ReadConfig("config.yml", conf); err != nil {
		panic(err)
	}
	conf.TemplateDir = "templates"
	if err := conf.EnsureFilesExist(); err != nil {
		panic(err)
	}
	storage := sqlite3.New(sqlite3.Config{Database: conf.SqliteDatabase})
	store := session.New(session.Config{
		CookieHTTPOnly: true,
		CookieSameSite: "true",
		Expiration:     35 * 24 * time.Hour,
		Storage:        storage})
	store.RegisterType("")
	store.RegisterType(1)
	handler := handlers.Handler{Conf: conf, Store: store}
	engine := html.New("./templates", ".html")
	engine.Reload(true)

	// Middlewares
	app := fiber.New(fiber.Config{Views: engine})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		logrus.Info("Gracefully shutting down...")
		_ = app.Shutdown()
	}()
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")
		return c.Next()
	})
	app.Use(recover.New())
	app.Use(logger.New())

	// Routes
	app.Static("/2fa-web", "./2fa-web")
	app.Get("/login", handler.LoginGetHandler)
	app.Post("/login", handler.LoginPostHandler)
	app.Post("/login/resend", handler.LoginResendHandler)
	app.Get("/login-check/:code", handler.LoginCheckHandler)
	app.Get("/logout", handler.LogoutGetHandler)
	app.Get("/password-reset", handler.PasswordResetGet)
	app.Post("/password-reset", handler.PasswordResetPost)
	app.Get("/password-change/:code", handler.PasswordChangeGet)
	app.Post("/password-change/:code", handler.PasswordChangePost)

	app.Get("/*", handler.SpaGetHandler)
	app.Listen(":3000")
}
