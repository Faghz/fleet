package util

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateInvoiceNumber(t *testing.T) {
	tests := []struct {
		name         string
		invoiceCount int64
		expected     string
	}{
		{
			name:         "First invoice for user",
			invoiceCount: 0,
			expected:     "INV-2025-000001",
		},
		{
			name:         "Fifth invoice for user",
			invoiceCount: 4,
			expected:     "INV-2025-000005",
		},
		{
			name:         "100th invoice for user",
			invoiceCount: 99,
			expected:     "INV-2025-000100",
		},
		{
			name:         "Large number invoice",
			invoiceCount: 123456,
			expected:     "INV-2025-123457",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoiceNumber := GenerateInvoiceNumber(tt.invoiceCount)

			assert.NotEmpty(t, invoiceNumber)

			// Check pattern: INV-YYYY-NNNNNN
			pattern := regexp.MustCompile(`^INV-\d{4}-\d{6}$`)
			assert.True(t, pattern.MatchString(invoiceNumber),
				"Invoice number should match pattern INV-YYYY-NNNNNN, got: %s", invoiceNumber)

			// Check year matches current year (2025)
			currentYear := time.Now().Year()
			expectedPrefix := fmt.Sprintf("INV-%d-", currentYear)
			assert.True(t, strings.HasPrefix(invoiceNumber, expectedPrefix),
				"Invoice number should start with %s, got: %s", expectedPrefix, invoiceNumber)

			// For the test year 2025, check the exact expected number
			if currentYear == 2025 {
				assert.Equal(t, tt.expected, invoiceNumber)
			}
		})
	}
}
