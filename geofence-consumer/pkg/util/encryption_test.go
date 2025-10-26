package util

import (
	"crypto/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name string
		text string
		salt string
	}{
		{
			name: "basic hash",
			text: "hello",
			salt: "world",
		},
		{
			name: "empty text",
			text: "",
			salt: "salt",
		},
		{
			name: "empty salt",
			text: "text",
			salt: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Hash(tt.text, tt.salt)
			assert.NotEmpty(t, result)

			// Test consistency - same inputs should produce same output
			result2 := Hash(tt.text, tt.salt)
			assert.Equal(t, result, result2)
		})
	}
}

func TestHash_Consistency(t *testing.T) {
	text := "test_password"
	salt := "test_salt"

	// Hash the same input multiple times
	hash1 := Hash(text, salt)
	hash2 := Hash(text, salt)

	// Should produce the same result
	assert.Equal(t, hash1, hash2, "Hash should be deterministic")
}

func TestEncryptData(t *testing.T) {
	// Generate a 32-byte key for AES-256
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := []byte("Hello, World!")
	saltLength := 16

	encrypted, err := EncryptData(plaintext, key, saltLength)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)

	// Verify the encrypted data is base64 encoded
	assert.True(t, isBase64(encrypted), "Encrypted data should be base64 encoded")
}

func TestEncryptData_EmptyPlaintext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := []byte("")
	saltLength := 16

	encrypted, err := EncryptData(plaintext, key, saltLength)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
}

func TestEncryptData_InvalidKey(t *testing.T) {
	// Invalid key length
	key := []byte("short")
	plaintext := []byte("test")
	saltLength := 16

	_, err := EncryptData(plaintext, key, saltLength)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create cipher")
}

func TestDecryptData(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	originalPlaintext := []byte("Hello, World!")
	saltLength := 16

	// Encrypt the data
	encrypted, err := EncryptData(originalPlaintext, key, saltLength)
	require.NoError(t, err)

	// Decrypt the data
	decrypted, err := DecryptData([]byte(encrypted), key, saltLength)
	require.NoError(t, err)

	assert.Equal(t, originalPlaintext, decrypted)
}

func TestDecryptData_InvalidCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	invalidCiphertext := []byte("invalid_base64!")
	saltLength := 16

	_, err = DecryptData(invalidCiphertext, key, saltLength)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode base64")
}

func TestDecryptData_TooShortCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Create a valid base64 string that's too short
	shortCiphertext := []byte("dGVzdA==") // "test" in base64
	saltLength := 16

	_, err = DecryptData(shortCiphertext, key, saltLength)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ciphertext too short")
}

func TestDecryptData_WrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	_, err := rand.Read(key1)
	require.NoError(t, err)
	_, err = rand.Read(key2)
	require.NoError(t, err)

	plaintext := []byte("secret message")
	saltLength := 16

	// Encrypt with key1
	encrypted, err := EncryptData(plaintext, key1, saltLength)
	require.NoError(t, err)

	// Try to decrypt with key2
	_, err = DecryptData([]byte(encrypted), key2, saltLength)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decrypt data")
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	testCases := []struct {
		name       string
		plaintext  string
		saltLength int
	}{
		{"short text", "hello", 8},
		{"medium text", "This is a longer message for testing", 16},
		{"long text", strings.Repeat("a", 1000), 32},
		{"empty text", "", 16},
		{"special characters", "!@#$%^&*()_+-=[]{}|;':\",./<>?", 16},
		{"unicode", "Hello 世界 🌍", 16},
	}

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plaintext := []byte(tc.plaintext)

			encrypted, err := EncryptData(plaintext, key, tc.saltLength)
			require.NoError(t, err)

			decrypted, err := DecryptData([]byte(encrypted), key, tc.saltLength)
			require.NoError(t, err)

			assert.Equal(t, plaintext, decrypted)
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "mypassword"
	salt := "mysalt"

	hash, err := HashPassword(password, salt)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Verify the hash is different each time (due to bcrypt salt)
	hash2, err := HashPassword(password, salt)
	require.NoError(t, err)
	assert.NotEqual(t, hash, hash2, "bcrypt should produce different hashes each time")
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	salt := "mysalt"

	_, err := HashPassword("", salt)
	assert.Error(t, err)
	assert.Equal(t, "password is empty", err.Error())
}

func TestHashPassword_EmptySalt(t *testing.T) {
	password := "mypassword"

	_, err := HashPassword(password, "")
	assert.Error(t, err)
	assert.Equal(t, "salt is empty", err.Error())
}

func TestComparePassword(t *testing.T) {
	password := "mypassword"
	salt := "mysalt"

	// Generate a hash
	hash, err := HashPassword(password, salt)
	require.NoError(t, err)

	// Compare with correct password
	err = ComparePassword(password, salt, string(hash))
	assert.NoError(t, err)

	// Compare with wrong password
	err = ComparePassword("wrongpassword", salt, string(hash))
	assert.Error(t, err)

	// Compare with wrong salt
	err = ComparePassword(password, "wrongsalt", string(hash))
	assert.Error(t, err)
}

func TestComparePassword_EmptyPassword(t *testing.T) {
	salt := "mysalt"
	hash := "somehash"

	err := ComparePassword("", salt, hash)
	assert.Error(t, err)
	assert.Equal(t, "password is empty", err.Error())
}

func TestComparePassword_EmptyHash(t *testing.T) {
	password := "mypassword"
	salt := "mysalt"

	err := ComparePassword(password, salt, "")
	assert.Error(t, err)
	assert.Equal(t, "hash is empty", err.Error())
}

func TestHashPassword_ComparePassword_Integration(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		salt     string
	}{
		{"simple", "password123", "salt123"},
		{"special chars", "p@ssw0rd!", "s@lt!"},
		{"unicode", "пароль", "соль"},
		{"medium password", strings.Repeat("a", 50), "salt"},
		{"long salt", "password", strings.Repeat("s", 30)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(tc.password, tc.salt)
			require.NoError(t, err)

			// Verify correct password passes
			err = ComparePassword(tc.password, tc.salt, string(hash))
			assert.NoError(t, err)

			// Verify wrong password fails
			err = ComparePassword(tc.password+"wrong", tc.salt, string(hash))
			assert.Error(t, err)
		})
	}
}

// Benchmark tests
func BenchmarkHash(b *testing.B) {
	text := "benchmark_text"
	salt := "benchmark_salt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hash(text, salt)
	}
}

func BenchmarkEncryptData(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	plaintext := []byte("benchmark plaintext data")
	saltLength := 16

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncryptData(plaintext, key, saltLength)
	}
}

func BenchmarkDecryptData(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	plaintext := []byte("benchmark plaintext data")
	saltLength := 16

	encrypted, _ := EncryptData(plaintext, key, saltLength)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DecryptData([]byte(encrypted), key, saltLength)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmark_password"
	salt := "benchmark_salt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashPassword(password, salt)
	}
}

// Helper function to check if a string is valid base64
func isBase64(s string) bool {
	// Check if string contains only valid base64 characters
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') ||
			c == '+' || c == '/' || c == '=') {
			return false
		}
	}
	return true
}
