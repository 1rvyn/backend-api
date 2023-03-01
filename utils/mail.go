package utils

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

func VerifyEmail(email string, code int) bool {
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

func SendMail(email string) error {
	// send a verification email

	fmt.Println("mailgun verification sending...")

	var mgDomain string = "api.irvyn.xyz"
	var mgApiKey string = os.Getenv("MAILGUN_API_KEY")
	mg := mailgun.NewMailgun(mgDomain, mgApiKey)

	mg.SetAPIBase("https://api.eu.mailgun.net/v3")

	// Build the email message
	from := "verifcation@irvyn.xyz"
	subject := "Account E-mail Verification"
	body := "Please confirm your email address by clicking the link below: \n\n https://api.irvyn.xyz/verify?code=" + strconv.Itoa(GenerateVerficiationCode()) + "&email=" + email

	message := mg.NewMessage(from, subject, body, email)

	// Send the email via the Mailgun API
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		fmt.Println("error sending email:")
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	}
	return nil
}

func GenerateVerficiationCode() int {

	rand.Seed(time.Now().UnixNano())

	code := rand.Intn(999999)

	fmt.Println("just generated a code: ", code)

	return code
}
