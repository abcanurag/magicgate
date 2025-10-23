package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
)

const (
	gcmNonceSize = 12 // 96 bits
)

func encrypt(key []byte, algoName string, plaintext []byte) ([]byte, error) {
	if !strings.EqualFold(algoName, "AES-256-GCM") {
		return nil, fmt.Errorf("unsupported algorithm: %s", algoName)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not create GCM cipher: %w", err)
	}

	// Nonce must be unique for each encryption with the same key
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Seal will append the ciphertext and authentication tag to the nonce
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(key []byte, algoName string, encryptedData []byte) ([]byte, error) {
	if !strings.EqualFold(algoName, "AES-256-GCM") {
		return nil, fmt.Errorf("unsupported algorithm: %s", algoName)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not create GCM cipher: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (data may be corrupt or key incorrect): %w", err)
	}

	return plaintext, nil
}
