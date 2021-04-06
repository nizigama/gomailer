package mailer

import (
	"encoding/json"
	"fmt"
	"gomailer/gmail"
	"gomailer/mailgun"
	"io/ioutil"
	"os"
)

type mailerConfig struct {
	DefaultSender string                `json:"sender"`
	Provider      string                `json:"provider"`
	MailgunConfig mailgun.MailgunConfig `json:"mailgun"`
	GmailConfig   gmail.GmailConfig     `json:"gmail"`
}

type Message struct {
	Subject string
	Body    string
}

var configs mailerConfig

const (
	MailgunProvider = "mailgun"
	GmailProvider   = "gmail"
)

func Init(workingDir string) error {

	file, err := os.Open(workingDir + "/config.json") // For read access.
	if err != nil {
		err = createConfigFile(workingDir)
		if err != nil {
			return err
		}
	} else {
		err = parseConfigFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func createConfigFile(workingDir string) error {
	configFile, err := os.Create(workingDir + "/config.json")

	if err != nil {
		return err
	}
	defer configFile.Close()

	configs := mailerConfig{
		Provider:      MailgunProvider,
		MailgunConfig: mailgun.MailgunConfig{},
		GmailConfig:   gmail.GmailConfig{},
	}

	bx, err := json.Marshal(configs)
	if err != nil {
		return err
	}

	jsonData := string(bx)

	fmt.Fprint(configFile, jsonData) // could also use file's write and sync methods
	return nil
}

func parseConfigFile(file *os.File) error {
	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &configs)
	if err != nil {
		return err
	}
	return nil
}

// Send sends the message with provided sender and recipient's email
func (m Message) Send(recipients ...string) error {
	status, id, err := mailgun.SendTextMessage(configs.MailgunConfig.Domain, configs.MailgunConfig.ApiKey, configs.DefaultSender, m.Subject, m.Body, recipients)
	fmt.Println(status, id)
	if err != nil {
		return err
	}

	return nil
}
