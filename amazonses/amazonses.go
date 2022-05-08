package amazonses

import (
	"context"
	"net/mail"

	"github.com/uhey22e/hedwig"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// Client for Amazon SES.
// Using SESv2 API with aws-sdk-go-v2.
type Mailer struct {
	client *sesv2.Client
}

func OpenMailer(ctx context.Context, cfg aws.Config, opts ...hedwig.MailerOptionFn) *hedwig.Mailer {
	d := &Mailer{
		client: sesv2.NewFromConfig(cfg),
	}
	return hedwig.NewMailer(d, opts...)
}

func (c *Mailer) SendMail(ctx context.Context, from mail.Address, to []mail.Address, m *hedwig.Mail) error {
	params := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(from.String()),
		Destination: &sesTypes.Destination{
			ToAddresses: make([]string, len(to)),
		},
		Content: &sesTypes.EmailContent{
			Simple: &sesTypes.Message{
				Subject: &sesTypes.Content{
					Data: &m.Subject,
				},
			},
		},
	}
	for i, t := range to {
		params.Destination.ToAddresses[i] = t.String()
	}
	switch m.ContentType {
	case hedwig.ContentTypePlainText:
		params.Content.Simple.Body = &sesTypes.Body{
			Text: &sesTypes.Content{
				Data: aws.String(m.String()),
			},
		}
	case hedwig.ContentTypeHTML:
		params.Content.Simple.Body = &sesTypes.Body{
			Html: &sesTypes.Content{
				Data: aws.String(m.String()),
			},
		}
	}
	_, err := c.client.SendEmail(ctx, params)
	return err
}
