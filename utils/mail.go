package utils

import (
	"net/smtp"
	"strconv"
)

func verifyEmail(email string, code int) bool {
	// Email verification
	from := "irvynhall@gmail.com"

	subject := "Verify your email"

	body := "Your verification code is: " + strconv.Itoa(code)

	msg := "From: " + from + "\n" + "Subject: " + subject + "\n" + body

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, "password", "smtp.gmail.com"), from, []string{email}, []byte(msg))
	if err != nil {
		return false
	} else {
		return true
	}

}
