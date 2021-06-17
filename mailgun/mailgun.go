package mailgun

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

func SendTextMessage(domain, apiKey, sender, subject, body string, recipients []string, isFromEU bool) (string, string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(sender, subject, body, recipients...)
	if isFromEU {
		mg.SetAPIBase(mailgun.APIBaseEU)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	statusMessage, messageID, err := mg.Send(ctx, m)
	return statusMessage, messageID, err
}
