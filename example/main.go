package main

import (
	"context"
	"html/template"
	"io"
	"net/mail"
	"net/smtp"

	"github.com/uhey22e/hedwig"
	"github.com/uhey22e/hedwig/generalsmtp"
)

func basic() {
	from := mail.Address{Address: "from@example.com"}
	auth := smtp.PlainAuth("", from.Address, "yourpassword", "localhost")
	client, _ := generalsmtp.OpenMailer(context.TODO(), "localhost:1025", auth, hedwig.DefaultFrom(from))
	to := []mail.Address{
		{Address: "to@example.com"},
	}
	msg := &hedwig.Mail{
		Subject: "Subject",
	}
	// hedwig.EMail has io.Writer interface to write the message body.
	io.WriteString(msg, "Hello world.")
	client.SendMail(context.TODO(), to, msg)
}

func withTemplate() {
	from := mail.Address{Address: "from@example.com"}
	client, _ := generalsmtp.OpenMailer(context.TODO(), "localhost:1025", nil, hedwig.DefaultFrom(from))
	to := []mail.Address{
		{Address: "to@example.com"},
	}
	msg := &hedwig.Mail{
		Subject:     "日本語を含む件名",
		ContentType: hedwig.ContentTypeHTML,
	}
	tmpl, _ := template.New("").Parse(`<p>Hello {{ . }}.</p><p>こんにちは、{{ . }}。</p>`)
	tmpl.Execute(msg, "Bob")
	client.SendMail(context.TODO(), to, msg)
}

func main() {
	basic()
	withTemplate()
}
