package hedwig

import (
	"context"

	"github.com/uhey22e/hedwig/driver"
	"github.com/uhey22e/hedwig/types"
)

type Mailer struct {
	d driver.Mailer
}

func NewMailer(d driver.Mailer) *Mailer {
	return &Mailer{d: d}
}

func (m *Mailer) SendMail(ctx context.Context, mail *types.Mail) error {
	return m.d.SendMail(ctx, mail)
}
