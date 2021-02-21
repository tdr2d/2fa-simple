package handlers

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func (handler *Handler) SpaGetHandler(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	if session.Get("email") == nil || session.Get("login_date_unix") == nil {
		return c.Redirect("/login")
	}

	path := filepath.Join(handler.Conf.SpaDirectory, filepath.Clean(c.Path()))
	info, err := os.Stat(path)
	fallback_path := filepath.Join(handler.Conf.SpaDirectory, handler.Conf.SpaFallback)

	if os.IsNotExist(err) || info.IsDir() {
		return c.SendFile(fallback_path)
	} else {
		return c.SendFile(path)
	}
}
