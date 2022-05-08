# hedwig

An email client for Golang, supports multiple email services.

You can use following sending methods.

- SMTP
- [Amazon SES](https://aws.amazon.com/ses/)

## Installation

This package can be installed with the go get command:

```sh
go get github.com/uhey22e/hedwig
```

## Usage

You can send an email via the `hedwig.Client` interface.
SMTP, Amazon SES and other clients implements this interface.

Basic usage - send an email with Gmail.

```go
import (
	"context"
	"net/mail"
	"net/smtp"

	"github.com/uhey22e/hedwig"
	"github.com/uhey22e/hedwig/generalsmtp"
)

from := mail.Address{Address: "from@example.com"}
auth := smtp.PlainAuth("", from.Address, "yourpassword", "smtp.gmail.com")
client, _ := generalsmtp.OpenMailer(context.TODO(), "smtp.gmail.com:587", auth, hedwig.DefaultFrom(from))
to := []mail.Address{
	{Address: "to@example.com"},
}
msg := &hedwig.Mail{
	Subject: "Subject",
}
// hedwig.EMail has io.Writer interface to write the message body.
io.WriteString(msg, "Hello world.")
client.SendMail(context.TODO(), from, to, msg)
```

Or you can use [html/template](https://pkg.go.dev/html/template) to render the message body.

```go
msg := &hedwig.Mail{
	Subject:     "Subject",
	ContentType: hedwig.ContentTypeHTML,
}
tmpl, _ := template.New("").Parse(`<p>Hello {{ . }}.</p>`)
tmpl.Execute(msg, "Bob")
client.SendMail(context.TODO(), from, to, msg)
```

## Supported services

### Amazon SES

You can use `github.com/uhey22e/hedwig/amazonses` driver.
This driver uses [aws-sdk-go-v2](github.com/aws/aws-sdk-go-v2) package.

```go
import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/uhey22e/hedwig"
	"github.com/uhey22e/hedwig/amazonses"
)

ctx := context.TODO()
cfg, _ := config.LoadDefaultConfig(ctx)
client := amazonses.OpenMailer(ctx, cfg)
```
