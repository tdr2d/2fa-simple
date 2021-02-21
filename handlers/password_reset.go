package handlers

import "github.com/gofiber/fiber/v2"

func (handler *Handler) PasswordResetGet(c *fiber.Ctx) error {
	return c.Render("password_reset", fiber.Map{"Title": "Password Reset"}, "layout")
}

type PasswordResetForm struct {
	Email string
}

type PasswordChangeForm struct {
	NewPassword      string
	NewPasswordCheck string
}

func (handler *Handler) PasswordResetPost(c *fiber.Ctx) error {
	// Step 1 send mail
	// session with email and check code
	return c.Render("password_reset", fiber.Map{"Title": "Password Reset"}, "layout")
}

func (handler *Handler) PasswordChangePost(c *fiber.Ctx) error {
	// Step 2 change password form
	// get mail in session
	return c.Render("password_reset", fiber.Map{"Title": "Password Reset"}, "layout")
}
