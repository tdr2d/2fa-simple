package handlers

import "github.com/gofiber/fiber/v2"

func (hand *Handler) ContentGetHandler(c *fiber.Ctx) error {
	return c.Render("content", fiber.Map{"Title": "Content"}, "layout")
}
