package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// GenerateAESKey generates a new AES-256 key (32 bytes)
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256 requires a 32-byte key
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %w", err)
	}
	return key, nil
}

// EncryptAESGCM encrypts plaintext using AES-256 GCM with a given key
// Returns base64 encoded ciphertext (nonce + ciphertext + tag)
func EncryptAESGCM(plaintext []byte, key []byte, nonceSize int) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	// Prepend nonce to ciphertext for easier decryption
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

// DecryptAESGCM decrypts base64 encoded ciphertext using AES-256 GCM with a given key
func DecryptAESGCM(base64Ciphertext string, key []byte, nonceSize int) ([]byte, error) {
	ciphertextWithNonce, err := base64.StdEncoding.DecodeString(base64Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	if len(ciphertextWithNonce) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short to contain nonce")
	}

	nonce := ciphertextWithNonce[:nonceSize]
	ciphertext := ciphertextWithNonce[nonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	return plaintext, nil
}