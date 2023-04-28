package utils

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func SendMail(email string, code int) error {
	// send a verification email

	fmt.Println("mailgun verification sending...")

	var mgDomain = "api.irvyn.xyz"
	var mgApiKey = os.Getenv("MAILGUN_API_KEY")
	mg := mailgun.NewMailgun(mgDomain, mgApiKey)

	mg.SetAPIBase("https://api.eu.mailgun.net/v3")

	// Build the email message
	from := "verifcation@irvyn.xyz"
	subject := "Account E-mail Verification"
	body := "Please confirm your email address by clicking the link below: \n\n https://api.irvyn.xyz/verify?code=" + strconv.Itoa(code) + "&email=" + email

	message := mg.NewMessage(from, subject, body, email)

	// Send the email via the Mailgun API
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		fmt.Println("error sending email:")
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
