package main

import (
	"os"
	"strings"
)

var DEFAULT_RECIPIENTS = []string{"kaneeldias@gmail.com"}
var DEFAULT_BCC = []string{"kaneeldias@gmail.com"}

func getRecipients() []string {
	recipients := os.Getenv("RECIPIENTS")
	if recipients == "" {
		return DEFAULT_RECIPIENTS
	}
	parts := strings.Split(recipients, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func getBCC() []string {
	bcc := os.Getenv("BCC")
	if bcc == "" {
		return DEFAULT_BCC
	}
	parts := strings.Split(bcc, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
