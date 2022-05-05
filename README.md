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

For example - send an email via the Gmail SMTP server.

```go
import (
	"context"
	"net/mail"
	"net/smtp"

	"github.com/uhey22e/hedwig"
)

var client hedwig.Client = &hedwig.SMTPClient{
	Address: "smtp.gmail.com:587",
	Auth:    smtp.PlainAuth("", "yourname@gmail.com", "password", "smtp.gmail.com"),
}

email := &hedwig.EMail{
	From: mail.Address{Address: "yourname@gmail.com"},
	To: []mail.Address{
		{Address: "to@gmail.com"},
	},
	Subject: "Subject",
	Body:    "Hello world.",
}
client.SendMail(context.TODO(), email)
```

### Amazon SES

You can use `AmazonSESClient`.

```go
import (
	"context"
	"net/mail"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/uhey22e/hedwig"
)

var client hedwig.Client

ctx := context.TODO()
cfg, _ := config.LoadDefaultConfig(ctx)
client = hedwig.NewAmazonSESClient(cfg)
```
