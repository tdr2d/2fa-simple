package handler

import (
	"2fa-simple/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

const check_mail_template = `You attempted to login to Ecorp.co<br>
Your verification code is <strong>%s</strong> (Valid 1h) <br>
<br>
Regards,<br>
<br>
<small>Ecorp.co</small>`

func LoginHandler(c *fiber.Ctx) error {
	var err error
	u := new(UserLogin)

	if err = c.BodyParser(u); err != nil {
		return err
	}

	code, err := utils.GenerateOTP(5)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if err = utils.SendGrid(
		utils.Recipient{Name: "Login Check", Mail: service_email}, // TODO parse config
		utils.Recipient{Name: "", Mail: u.Email},
		"Login Mail check",
		fmt.Sprintf(check_mail_template, code)); err != nil {
		logrus.Error(err)
		return err
	}

	return c.SendString("Hello, World!")
}
