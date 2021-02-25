package handlers

import (
	"2fa-simple/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Conf  *utils.Config
	Store *session.Store
}

type UserLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (handler *Handler) sendLoginCheckMail(session *session.Session, email string) error {
	code, err := utils.GenerateOTP(15)
	if err != nil {
		return err
	}
	session.Set("login_code", code)
	session.Set("login_code_expiration", int(time.Now().Add(time.Hour).Unix()))
	logrus.Info(session.Get("login_code_expiration"))
	data := map[string]string{
		"Website": handler.Conf.Website,
		"Link":    fmt.Sprintf("%s/login-check/%s", handler.Conf.BaseUrl, code),
		"Company": handler.Conf.CompanyName,
	}
	login_check_content := utils.TranslateWithArgs(handler.Conf.Language, "mail_login_check", data)
	// logrus.Info(login_check_content)
	if err = utils.SendGrid(
		utils.MailUser{Name: utils.Translate(handler.Conf.Language, "mail_login_check_name"), Mail: handler.Conf.ServiceEmail},
		utils.MailUser{Name: "", Mail: email},
		utils.Translate(handler.Conf.Language, "mail_login_check_title"),
		login_check_content); err != nil {
		return err
	}
	return nil
}

func (handler *Handler) LoginPostHandler(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	u := new(UserLogin)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	configPasswordHash, err := handler.Conf.GetPasswordHashFromUserEmail(strings.TrimSpace(u.Email))
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(u.Password, configPasswordHash) {
		return c.Status(fiber.StatusUnauthorized).SendString(utils.Translate(handler.Conf.Language, "user_password_mismatch"))
	}
	session.Set("email", strings.TrimSpace(u.Email))
	if err := handler.sendLoginCheckMail(session, strings.TrimSpace(u.Email)); err != nil {
		return err
	}
	session.Save()
	return c.SendString("ok")
}

func (handler *Handler) LoginResendHandler(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	if session.Get("email") == nil || session.Get("login_code") == nil || session.Get("login_code_expiration") == nil {
		return c.Status(fiber.StatusUnauthorized).SendString(utils.Translate(handler.Conf.Language, "email_code_code_expiration_undefined"))
	}
	if err := handler.sendLoginCheckMail(session, strings.TrimSpace(session.Get("email").(string))); err != nil {
		return err
	}
	session.Save()
	return c.SendString("ok")
}

func (handler *Handler) LoginCheckHandler(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	if session.Get("email") == nil || session.Get("login_code") == nil || session.Get("login_code_expiration") == nil {
		return c.Status(fiber.StatusUnauthorized).SendString(utils.Translate(handler.Conf.Language, "email_code_code_expiration_undefined"))
	}

	if int(time.Now().Unix()) > session.Get("login_code_expiration").(int) {
		return c.Status(fiber.StatusUnauthorized).SendString(utils.Translate(handler.Conf.Language, "expired_link"))
	}

	if c.Params("code") != session.Get("login_code").(string) {
		return c.Status(fiber.StatusUnauthorized).SendString(utils.Translate(handler.Conf.Language, "invalid_link"))
	}

	session.Set("login_date_unix", int(time.Now().Unix()))
	session.Delete("login_code")
	session.Delete("login_code_expiration")
	session.Save()
	return c.Redirect("/")
}

func (handler *Handler) LoginGetHandler(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{"Title": "Login", "lang": handler.Conf.Language}, "layout")
}

func (handler *Handler) LogoutGetHandler(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	session.Destroy()
	return c.Redirect("/login")
}
