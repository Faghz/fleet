package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func GeneratePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?"

	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}
	return string(password), nil
}

// GenerateInvoiceNumber generates a sequential invoice number in format INV-YYYY-NNNNNN
// The NNNNNN part is incremental based on the current invoice count + 1
func GenerateInvoiceNumber(currentCount int64) string {
	// Get current year
	year := time.Now().Year()

	// Increment count by 1 for the new invoice
	nextNumber := currentCount + 1

	// Format as INV-YYYY-NNNNNN (zero-padded to 6 digits)
	invoiceNumber := fmt.Sprintf("INV-%d-%06d", year, nextNumber)

	return invoiceNumber
}
