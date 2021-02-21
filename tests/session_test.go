package main

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"github.com/gofiber/utils"
	"github.com/valyala/fasthttp"
)

func Test_Session(t *testing.T) {
	storage := sqlite3.New()
	store := session.New(session.Config{
		Storage: storage,
	})
	store.RegisterType("")
	store.RegisterType(1)

	// fiber instance and context
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// set cookie
	ctx.Request().Header.SetCookie(store.CookieName, "123")

	// get session
	sess, _ := store.Get(ctx)
	sess.Set("name", "john")
	sess.Set("code", 1234)
	sess.Save()

	// get value after save
	sess, err := store.Get(ctx)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "john", sess.Get("name").(string))
	utils.AssertEqual(t, 1234, sess.Get("code").(int))
}
