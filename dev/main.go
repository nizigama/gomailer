package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

const (
	apiKey string = "xxx"
	domain string = "xxx@domain.com"
)

func main() {
	receivers := []string{
		"first@example.com", "second@example.com", "third@example.com",
	}

	mg := mailgun.NewMailgun(domain, apiKey)

	mg.SetAPIBase(mailgun.APIBaseEU)

	newMessage := mg.NewMessage("test@example.com", "Testing multiple receivers", "Alloooooooooooooooooooooooooooo", receivers...)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	statusMessage, messageID, err := mg.Send(ctx, newMessage)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Everything went well")
	fmt.Println(statusMessage)
	fmt.Println(messageID)
}
