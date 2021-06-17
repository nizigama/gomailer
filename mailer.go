package gomailer

import (
	"errors"
	"fmt"

	"github.com/nizigama/gomailer/mailgun"
)

// Message type contains the sender's email, the email's subject and body
// PS: the sender's email can be left empty but make sure to initialize the default sender
// with the SetDefaultSender function
type Message struct {
	Sender  string
	Subject string
	Body    string
}

type mailgunCredentials struct {
	domain string
	apiKey string
}

var credentials mailgunCredentials
var defaultSender string

// SetCredentials takes the mailgun domain and the api key to initiate the connection with mailgun servers
func SetCredentials(mailgunDomain, apiKey string) error {

	if len(mailgunDomain) == 0 || len(apiKey) == 0 {
		return errors.New("invalid credentials")
	}

	credentials = mailgunCredentials{
		domain: mailgunDomain,
		apiKey: apiKey,
	}

	return nil
}

// SetDefaultSender sets the default sender's email which helps in case you want to send multiple messages
// without always specifying the sender
func SetDefaultSender(senderEmail string) {
	defaultSender = senderEmail
}

// Send sends the message with provided sender and recipient's email
func (m Message) Send(recipients ...string) error {

	var messageSender string

	if m.Sender == "" {
		messageSender = defaultSender
	} else {
		messageSender = m.Sender
	}

	status, id, err := mailgun.SendTextMessage(credentials.domain, credentials.apiKey, messageSender, m.Subject, m.Body, recipients)
	fmt.Println(status, id)
	if err != nil {
		return err
	}

	return nil
}
