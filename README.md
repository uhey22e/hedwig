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

```go
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

### SMTP

You can use `SMTPClient`.

### Amazon SES

You can use `AmazonSESClient`.

```go
import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
)

ctx := context.TODO()
cfg, _ := config.LoadDefaultConfig(ctx)
var client hedwig.Client = NewAmazonSESClient(cfg)
```
