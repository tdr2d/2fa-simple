package handlers

import (
	"2fa-simple/utils"
	"fmt"
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
	data := fiber.Map{
		"Website": "mywebsite",
		"Link":    fmt.Sprintf("%s/login-check/%s", handler.Conf.BaseUrl, code),
		"Company": handler.Conf.CompanyName,
	}
	login_check_content, err := utils.RenderTemplate(fmt.Sprintf("%s/mail_login_check.html", handler.Conf.TemplateDir), data)
	if err != nil {
		return err
	}
	logrus.Info(login_check_content)
	// if err = utils.SendGrid(
	// 	utils.MailUser{Name: "Login Check", Mail: handler.Conf.ServiceEmail},
	// 	utils.MailUser{Name: "", Mail: email},
	// 	"Login Mail check",
	// 	login_check_content); err != nil {
	// 	return err
	// }
	return nil
}

func (handler *Handler) LoginPostHandler(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	u := new(UserLogin)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	configPasswordHash, err := handler.Conf.GetPasswordHashFromUserEmail(u.Email)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(u.Password, configPasswordHash) {
		c.SendStatus(fiber.StatusUnauthorized)
		return c.SendString("user_password_mismatch")
	}
	session.Set("email", u.Email)
	if err := handler.sendLoginCheckMail(session, u.Email); err != nil {
		return err
	}
	session.Save()
	return c.SendString("ok")
}

func (handler *Handler) LoginResendHandler(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	if session.Get("email") == nil || session.Get("login_code") == nil || session.Get("login_code_expiration") == nil {
		c.SendStatus(fiber.StatusUnauthorized)
		return c.SendString("email_code_code_expiration_undefined")
	}
	if err := handler.sendLoginCheckMail(session, session.Get("email").(string)); err != nil {
		return err
	}
	session.Save()
	return c.SendString("ok")
}

func (handler *Handler) LoginCheckHandler(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	if session.Get("email") == nil || session.Get("login_code") == nil || session.Get("login_code_expiration") == nil {
		c.SendStatus(fiber.StatusUnauthorized)
		return c.SendString("email_code_code_expiration_undefined")
	}

	if int(time.Now().Unix()) > session.Get("login_code_expiration").(int) {
		c.SendStatus(fiber.StatusUnauthorized)
		return c.SendString("code_expired")
	}

	code_param := c.Params("code")
	if code_param != session.Get("login_code").(string) {
		c.SendStatus(fiber.StatusUnauthorized)
		return c.SendString("code_invalid")
	}

	session.Set("login_date_unix", int(time.Now().Unix()))
	session.Delete("login_code")
	session.Delete("login_code_expiration")
	return c.Redirect("/")
}

func (handler *Handler) LoginGetHandler(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{"Title": "Login"}, "layout")
}

func (handler *Handler) LogoutGetHandler(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	session.Destroy()
	return c.Redirect("/login")
}
