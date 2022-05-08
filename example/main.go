package main

import (
	"context"
	"log"
	"net/mail"

	"github.com/uhey22e/hedwig/generalsmtp"
	"github.com/uhey22e/hedwig/types"
)

func main() {
	client, err := generalsmtp.NewMailer(context.TODO(), "localhost:1025", nil)
	if err != nil {
		log.Fatal(err)
	}
	email := &types.Mail{
		From: mail.Address{Address: "yourname@gmail.com"},
		To: []mail.Address{
			{Address: "to@gmail.com"},
		},
		Subject: "Subject",
		Body:    "Hello world.",
	}
	err = client.SendMail(context.TODO(), email)
	if err != nil {
		log.Fatal(err)
	}
}
