package main

import (
	"fmt"
	"log"

	"github.com/sethvargo/go-password/password"
	gomail "gopkg.in/mail.v2"
)

func main() {

}

func readConfig() {

}

func sendEmail(receiver string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "from@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", receiver)

	// Set E-Mail subject
	m.SetHeader("Subject", "Verification code.")

	// Generate verification code
	code, err := password.Generate(7, 7, 0, false, true)
	if err != nil {
		log.Fatal(err)
	}

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", code)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "from@gmail.com", "<email_password>")

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return
}

//func generateCode() string {

//}
