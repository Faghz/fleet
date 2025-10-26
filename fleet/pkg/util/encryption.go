package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func Hash(text, salt string) string {
	text = salt + text + salt
	hash := sha256.Sum256([]byte(text))
	return base64.StdEncoding.EncodeToString(hash[:]) // Store in Base64
}

// EncryptData encrypts the given plaintext by first adding a random salt
// to both the beginning and the end, and then encrypting the salted data using AES-GCM.
// The returned byte slice contains the nonce prepended to the ciphertext,
// which is necessary for decryption.
func EncryptData(plaintext, key []byte, saltLength int) (string, error) {
	// Define the length of the salt.

	// Generate a random salt.
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Append salt at the beginning and the end of the plaintext.
	saltedData := append(salt, plaintext...)
	saltedData = append(saltedData, salt...)

	// Create a new AES cipher using the provided key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Wrap the block cipher in GCM mode.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce.
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the salted data.
	ciphertext := aesGCM.Seal(nil, nonce, saltedData, nil)
	decryptedText := append(nonce, ciphertext...)

	// Prepend the nonce to the ciphertext.
	// (This is necessary for decryption; the nonce does not need to be secret.)
	return base64.StdEncoding.EncodeToString(decryptedText), nil
}

// DecryptData decrypts the ciphertext produced by EncryptData.
// It extracts the nonce, decrypts the data, and then removes the salt
// from the beginning and the end to recover the original plaintext.
func DecryptData(ciphertext, key []byte, saltLength int) ([]byte, error) {
	decodedStr, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create a new AES cipher with the given key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Wrap the cipher in GCM mode.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(decodedStr) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract the nonce from the start of the ciphertext.
	nonce, ct := decodedStr[:nonceSize], decodedStr[nonceSize:]

	// Decrypt the data.
	decrypted, err := aesGCM.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Ensure the decrypted data is long enough to contain both salts.
	if len(decrypted) < 2*saltLength {
		return nil, fmt.Errorf("decrypted data too short to contain salts")
	}

	if !bytes.Equal(decrypted[:saltLength], decrypted[len(decrypted)-saltLength:]) {
		return nil, fmt.Errorf("salt mismatch")
	}

	// Remove the salt from the beginning and end to get the original plaintext.
	plaintext := decrypted[saltLength : len(decrypted)-saltLength]
	return plaintext, nil
}

func HashPassword(pw, salt string) (res []byte, err error) {
	if salt == "" {
		return nil, errors.New("salt is empty")
	}
	if pw == "" {
		return nil, errors.New("password is empty")
	}

	password := pw + salt

	res, err = bcrypt.GenerateFromPassword([]byte(password), 13)

	return
}

func ComparePassword(pw, salt, hash string) (err error) {
	if pw == "" {
		return errors.New("password is empty")
	}
	if hash == "" {
		return errors.New("hash is empty")
	}

	pw = pw + salt

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return
}
