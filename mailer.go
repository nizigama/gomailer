package gomailer

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// Message type contains the sender's email, the email's subject and body
// PS: the sender's email can be left empty but make sure to initialize the default sender
// with the SetDefaultSender function or during package initialization with Init() function
type Message struct {
	Sender  string
	Subject string
	Body    string
}

type mailgunSettings struct {
	domain        string
	apiKey        string
	inEURegion    bool
	defaultSender string
}

var credentials mailgunSettings
var mg *mailgun.MailgunImpl

// Init initiates the mailgun configurations to start sending emails
func Init(mailgunDomain, apiKey, defaultSender string, isInEURegion bool) error {
	if len(strings.Trim(mailgunDomain, " ")) == 0 || len(strings.Trim(apiKey, " ")) == 0 {
		return errors.New("invalid credentials")
	}

	defaultSender = strings.Trim(defaultSender, " ")

	if defaultSender != "" {

		if err := ValidateEmail(defaultSender); err != nil {
			return err
		}
	}

	credentials = mailgunSettings{
		domain:        mailgunDomain,
		apiKey:        apiKey,
		inEURegion:    isInEURegion,
		defaultSender: defaultSender,
	}

	mg = mailgun.NewMailgun(credentials.domain, credentials.apiKey)

	if credentials.inEURegion {
		mg.SetAPIBase(mailgun.APIBaseEU)
	}

	return nil
}

// ValidateEmail verifies if the email has a valid email format
func ValidateEmail(email string) error {

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email")
	}

	parts := strings.Split(email, "@")

	// if there is no content before or after the @ symbol
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return fmt.Errorf("invalid email")
	}

	afterAtSymbol := strings.Split(parts[1], ".")
	// if there is a dot after the @ symbol
	if len(afterAtSymbol) < 2 {
		return fmt.Errorf("invalid email")
	}

	// if there is content after the last dot(.)
	if strings.Trim(afterAtSymbol[len(afterAtSymbol)-1], " ") == "" {
		return fmt.Errorf("invalid email")
	}

	return nil
}

// SetDefaultSender sets the default sender's email which helps in case you want to send multiple messages
// without always specifying the sender.
// Returns an error if the email doesn't have a valid email format
func SetDefaultSender(senderEmail string) error {
	senderEmail = strings.Trim(senderEmail, " ")

	if err := ValidateEmail(senderEmail); err != nil {
		return err
	}

	credentials.defaultSender = senderEmail

	return nil
}

// Send sends a simple text email to the provided recipients' emails
func (m Message) SendSimpleTextEmail(recipients ...string) (string, string, error) {

	var messageSender string

	if strings.Trim(m.Sender, " ") == "" {

		if credentials.defaultSender != "" {
			messageSender = credentials.defaultSender
		} else {
			return "", "", fmt.Errorf("no default sender set")
		}

	} else {

		if err := ValidateEmail(m.Sender); err != nil {
			return "", "", err
		}

		messageSender = m.Sender
	}

	newMessage := mg.NewMessage(messageSender, m.Subject, m.Body, recipients...)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	statusMessage, messageID, err := mg.Send(ctx, newMessage)

	if err != nil {
		return "", "", nil
	}

	return statusMessage, messageID, nil
}

// Send sends an email with one attachment the provided recipients' emails
func (m Message) SendEmailWithOneFileAttachment(attachment *os.File, recipients ...string) (string, string, error) {

	var messageSender string

	if strings.Trim(m.Sender, " ") == "" {

		if credentials.defaultSender != "" {
			messageSender = credentials.defaultSender
		} else {
			return "", "", fmt.Errorf("no default sender set")
		}

	} else {
		messageSender = m.Sender
	}

	newMessage := mg.NewMessage(messageSender, m.Subject, m.Body, recipients...)

	fileBytes, err := ioutil.ReadAll(attachment)

	if err != nil {
		return "", "", err
	}

	newMessage.AddBufferAttachment(attachment.Name(), fileBytes)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	statusMessage, messageID, err := mg.Send(ctx, newMessage)

	if err != nil {
		return "", "", nil
	}

	return statusMessage, messageID, nil
}
