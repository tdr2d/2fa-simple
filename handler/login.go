package handler

import (
	"2fa-simple/utils"
	"fmt"
	"strconv"
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

type UserLogin2 struct {
	Code string `json:"code" form:"code"`
}

const check_mail_template = ``

func (hand *Handler) sendLoginCheckMail(session *session.Session, email string) error {
	code, err := utils.GenerateOTP(8)
	if err != nil {
		return err
	}
	session.Set("login_code", code)
	session.Set("login_code_expiration", int(time.Now().Add(time.Hour).Unix()))

	data := fiber.Map{
		"Website": "mywebsite",
		"Link":    fmt.Sprintf("%s/login-check/%s", hand.Conf.BaseUrl, code),
		"Company": hand.Conf.CompanyName,
	}
	login_check_content, err := utils.RenderTemplate("mail_login_check.html", data)
	if err != nil {
		return err
	}
	logrus.Info(login_check_content)
	// if err = utils.SendGrid(
	// 	utils.MailUser{Name: "Login Check", Mail: hand.Conf.ServiceEmail},
	// 	utils.MailUser{Name: "", Mail: email},
	// 	"Login Mail check",
	// 	login_check_content); err != nil {
	// 	return err
	// }
	return nil
}

func (hand *Handler) LoginPostHandler(c *fiber.Ctx) error {
	session := hand.GetSession(c)
	u := new(UserLogin)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	configPasswordHash, err := hand.Conf.GetPasswordHashFromUserEmail(u.Email)
	logrus.Info(configPasswordHash)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(u.Password, configPasswordHash) {
		c.SendStatus(400)
		return c.SendString("user_password_mismatch")
	}

	session.Set("email", u.Email)
	if err := hand.sendLoginCheckMail(session, u.Email); err != nil {
		logrus.Info(err)
		return err
	}

	session.Save()
	return c.SendString("ok")
}

func (hand *Handler) LoginResendHandler(c *fiber.Ctx) error {
	session := hand.GetSession(c)
	email := fmt.Sprintf("%v", session.Get("email"))
	code := fmt.Sprintf("%v", session.Get("login_code"))
	code_expiration := fmt.Sprintf("%v", session.Get("login_code_expiration"))
	if email == "" || code == "" || code_expiration == "" {
		c.SendStatus(400)
		return c.SendString("email_undefined")
	}
	if err := hand.sendLoginCheckMail(session, email); err != nil {
		return err
	}
	return c.SendString("ok")
}

func (hand *Handler) LoginCheckHandler(c *fiber.Ctx) error {
	session := hand.GetSession(c)
	email := fmt.Sprintf("%v", session.Get("email"))
	code := fmt.Sprintf("%v", session.Get("login_code"))
	code_expiration := fmt.Sprintf("%v", session.Get("login_code_expiration"))

	if email == "" || code == "" || code_expiration == "" {
		c.SendStatus(400)
		return c.SendString("email_code_code_expiration_undefined")
	}

	code_expiration_int, err := strconv.Atoi(code_expiration)
	if err != nil {
		return err
	}
	if code_expiration_int > int(time.Now().Unix()) {
		c.SendStatus(400)
		return c.SendString("code_expired") // TODO redirect to login screen
	}

	code_param := c.Params("code")
	if code_param != code {
		c.SendStatus(400)
		return c.SendString("code_invalid") // TODO redirect to login screen
	}

	session.Set("login_date_unix", int(time.Now().Unix()))
	session.Delete("login_code")
	session.Delete("login_code_expiration")
	return c.SendString("ok")
}

func (hand *Handler) LoginGetHandler(c *fiber.Ctx) error {
	error_message := c.Query("error_message", "")
	return c.Render("login", fiber.Map{"Title": "Login", "ErrorMsg": error_message}, "layout")
}
