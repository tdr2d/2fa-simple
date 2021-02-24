package handlers

import (
	"2fa-simple/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (handler *Handler) PasswordResetGet(c *fiber.Ctx) error {
	return c.Render("password_reset", fiber.Map{"Title": "Password Reset", "lang": handler.Conf.Language}, "layout")
}

type PasswordResetForm struct {
	Email string `json:"email" form:"email"`
}

// Step 1 send link in mail, save code in session
func (handler *Handler) PasswordResetPost(c *fiber.Ctx) error {
	u := new(PasswordResetForm)
	if err := c.BodyParser(u); err != nil {
		return err
	}
	if !handler.Conf.UserExists(strings.TrimSpace(u.Email)) {
		logrus.Info("User does not exist " + u.Email)
		return c.SendString("ok")
	}
	code, err := utils.GenerateOTP(15)
	if err != nil {
		return err
	}
	session := handler.GetSession(c)
	session.Set("email", strings.TrimSpace(u.Email))
	session.Set("password_change_code", code)
	session.Set("password_change_code_expiration", int(time.Now().Add(time.Hour).Unix()))
	data := fiber.Map{
		"Website": handler.Conf.Website,
		"Link":    fmt.Sprintf("%s/password-change/%s", handler.Conf.BaseUrl, code),
		"Company": handler.Conf.CompanyName,
		"lang":    handler.Conf.Language,
	}
	password_reset_mail, err := utils.RenderTemplate(handler.Conf.TemplateDir+"/mail/change_password.html", data)
	if err != nil {
		return err
	}

	logrus.Info(password_reset_mail)
	// if err = utils.SendGrid(
	// 	utils.MailUser{Name: "Password Reset", Mail: handler.Conf.ServiceEmail},
	// 	utils.MailUser{Name: "", Mail: email},
	// 	"Password Reset",
	// 	password_reset_mail); err != nil {
	// 	return err
	// }
	session.Save()
	return c.SendString("ok")
}

// Step 2 acccess link in mail
func (handler *Handler) PasswordChangeGet(c *fiber.Ctx) error {
	session := handler.GetSession(c)
	if session.Get("email") == nil || session.Get("password_change_code") == nil || session.Get("password_change_code_expiration") == nil {
		return fiber.ErrUnauthorized
	}
	if int(time.Now().Unix()) > session.Get("password_change_code_expiration").(int) {
		return c.Status(fiber.StatusUnauthorized).SendString("code_expired")
	}

	if c.Params("code") != session.Get("password_change_code").(string) {
		return c.Status(fiber.StatusUnauthorized).SendString("code_invalid")
	}

	return c.Render("password_change", fiber.Map{"Title": "Password Change", "lang": handler.Conf.Language}, "layout")
}

type PasswordChangeForm struct {
	NewPassword      string `json:"newpassword" form:"newpassword"`
	NewPasswordCheck string `json:"newpasswordcheck" form:"newpasswordcheck"`
}

// Step 3 change password form
func (handler *Handler) PasswordChangePost(c *fiber.Ctx) error {
	form := new(PasswordChangeForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	logrus.Info(form.NewPassword)
	logrus.Info(form.NewPasswordCheck)
	if form.NewPassword != form.NewPasswordCheck {
		return c.Status(fiber.StatusBadRequest).SendString("Passwords are not identical")
	}

	session := handler.GetSession(c)
	if session.Get("email") == nil || session.Get("password_change_code") == nil || session.Get("password_change_code_expiration") == nil {
		return fiber.ErrUnauthorized
	}
	if int(time.Now().Unix()) > session.Get("password_change_code_expiration").(int) {
		return c.Status(fiber.StatusUnauthorized).SendString("code_expired")
	}

	email := strings.TrimSpace(session.Get("email").(string))
	hash := utils.HashPassword(form.NewPassword)
	if err := handler.Conf.ChangePassword(email, hash); err != nil {
		return err
	}
	logrus.Info("New Password " + email + " " + hash)
	utils.WriteYaml(&handler.Conf, "config.yml")

	// TODO write yaml yaml
	session.Set("login_date_unix", int(time.Now().Unix()))
	session.Delete("password_change_code")
	session.Delete("password_change_code_expiration")
	session.Save()
	return c.SendString("ok")
}
