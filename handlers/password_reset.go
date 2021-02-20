package handlers

import "github.com/gofiber/fiber/v2"

func (handler *Handler) PasswordResetGet(c *fiber.Ctx) error {
	return c.Render("password_reset", fiber.Map{"Title": "Password Reset"}, "layout")
}

type PasswordResetForm struct {
	Email string
}

func (handler *Handler) PasswordResetPost(c *fiber.Ctx) error {
	return c.Render("password_reset", fiber.Map{"Title": "Password Reset"}, "layout")
}

func (handler *Handler) PasswordChangePost(c *fiber.Ctx) error {
	return c.Render("password_reset", fiber.Map{"Title": "Password Reset"}, "layout")
}
