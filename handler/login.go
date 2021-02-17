package handler

import (
	"2fa-simple/utils"
	"fmt"

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

type UserLogin2 struct {
	Code string `json:"code" form:"code"`
}

type ErrorUserPasswordNoMatch struct{}

func (e *ErrorUserPasswordNoMatch) Error() string {
	return "ErrorUserPasswordNoMatch"
}

const check_mail_template = `You attempted to login to Ecorp.co<br>
Your verification code is <strong>%s</strong> (Valid 1h) <br>
<br>
Regards,<br>
<br>
<small>Ecorp.co</small>`

func (hand *Handler) GetSession(c *fiber.Ctx) *session.Session {
	session, err := hand.Store.Get(c)
	if err != nil {
		panic(err)
	}
	return session
}

func (hand *Handler) LoginPostHandler(c *fiber.Ctx) error {
	session := hand.GetSession(c)
	u := new(UserLogin)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	configPasswordHash, err := hand.Conf.GetPasswordHashFromUserEmail(u.Email)
	if err != nil {
		return err
	}

	tmpPassword := utils.HashPassword(u.Password)
	logrus.Info(tmpPassword)
	if configPasswordHash != tmpPassword {
		return new(ErrorUserPasswordNoMatch)
	}

	code, err := utils.GenerateOTP(5)
	if err != nil {
		logrus.Error(err)
		return err
	}

	session.Set("email", u.Email)
	session.Set("login_code", code)
	if err = utils.SendGrid(
		utils.MailUser{Name: "Login Check", Mail: hand.Conf.ServiceEmail},
		utils.MailUser{Name: "", Mail: u.Email},
		"Login Mail check",
		fmt.Sprintf(check_mail_template, code)); err != nil {
		logrus.Error(err)
		return err
	}
	session.Save()
	return c.SendString("ok")
}

func (hand *Handler) LoginPostCheckHandler(c *fiber.Ctx) error {
	session := hand.GetSession(c)
	email := session.Get("email")
	code := session.Get("code")
	// TODO
}

func (hand *Handler) LoginGetHandler(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{"Title": "Login"}, "layout")
}
