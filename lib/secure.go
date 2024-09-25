// secure.go
package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var key = []byte("") // Use a fixed key for simplicity

// Encrypt the plaintext string
func encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12) // 12 bytes for GCM
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nil, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(nonce) + ":" + hex.EncodeToString(ciphertext), nil
}

// Decrypt the encrypted string
func decrypt(encrypted string) (string, error) {
	parts := []string{encrypted[:24], encrypted[25:]} // Split nonce and ciphertext
	nonce, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	ciphertext, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Write a string to a file securely
func SecureWrite(content string) error {
	encryptedString, err := encrypt(content)

	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}
	filePath := filepath.Join(homeDir, "github.sec") // Store in home directory
	fmt.Println("PATH IS:", filePath)

	return os.WriteFile(filePath, []byte(encryptedString), 0600) // Secure permissions
}

// Read a string from a file securely
func SecureRead() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(homeDir, "github.sec") // Get full path
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return decrypt(string(data))
}
