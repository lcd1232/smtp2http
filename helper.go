package main

import (
	"fmt"
	"net/mail"
)

// extractEmails returns emails.
func extractEmails(addresses []*mail.Address, _ ...error) []string {
	emails := make([]string, 0, len(addresses))

	for _, address := range addresses {
		emails = append(emails, fmt.Sprint(address))
	}

	return emails
}
