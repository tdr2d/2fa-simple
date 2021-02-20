package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func (hand *Handler) GetSession(c *fiber.Ctx) *session.Session {
	session, err := hand.Store.Get(c)
	if err != nil {
		panic(err)
	}
	return session
}
