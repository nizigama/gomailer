package helpers

import (
	"errors"
	"strings"
)

// VerifyMessageID verifies if a messageID contains certain symbols and characters that most messageIDs contain
func VerifyMessageID(messageID string) error {

	if !strings.Contains(messageID, "<") || !strings.Contains(messageID, ">") || !strings.Contains(messageID, "@") || !strings.Contains(messageID, ".") {
		return errors.New("invalid messageID/references format")
	}

	return nil
}
