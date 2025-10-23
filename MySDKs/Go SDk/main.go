package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	fmt.Println("--- Go SDK Usage Example ---\n")

	// 1. Initialize the SDK
	fmt.Println("1. Initializing SDK...")
	if err := Init("my-app-registration-token-12345"); err != nil {
		log.Fatalf("Failed to initialize SDK: %v", err)
	}
	fmt.Println()

	// 2. Create a session
	fmt.Println("2. Creating a session...")
	jwt, err := CreateSession("app_user_01", "super_secret_password")
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	fmt.Printf("   Session created. Received JWT: %.30s...\n\n", jwt)

	// 3. Manage a key (CREATE)
	fmt.Println("3. Creating a key named 'MySecretKey'...")
	keyMaterial := []byte("this-is-my-super-secret-key-data")
	if err := KeyOperation("CREATE", "MySecretKey", keyMaterial); err != nil {
		log.Fatalf("Failed to create key: %v", err)
	}
	fmt.Println("   Key 'MySecretKey' created on backend.\n")

	// 4. Perform crypto operations (Encrypt and Decrypt)
	fmt.Println("4. Encrypting and decrypting data with 'MySecretKey'...")
	plaintext := []byte("This is a very sensitive message.")
	algo := "AES-256-GCM"

	// Encrypt
	ciphertext, err := DoCrypto("ENCRYPT", "MySecretKey", algo, plaintext)
	if err != nil {
		log.Fatalf("Encryption failed: %v", err)
	}
	fmt.Println("   Encryption successful.")
	fmt.Printf("   Plaintext: '%s'\n", string(plaintext))
	fmt.Printf("   Ciphertext (hex): %s\n\n", hex.EncodeToString(ciphertext))

	// Decrypt
	decryptedTextBytes, err := DoCrypto("DECRYPT", "MySecretKey", algo, ciphertext)
	if err != nil {
		log.Fatalf("Decryption failed: %v", err)
	}
	fmt.Println("   Decryption successful.")
	fmt.Printf("   Decrypted Text: '%s'\n\n", string(decryptedTextBytes))

	// Verify
	if !bytes.Equal(plaintext, decryptedTextBytes) {
		log.Fatal("FATAL: Decrypted text does not match original plaintext!")
	}
	fmt.Println("   SUCCESS: Original plaintext matches decrypted text.\n")

	// 5. Demonstrate key caching
	fmt.Println("5. Encrypting again (should use cached key)...")
	ciphertext2, err := DoCrypto("ENCRYPT", "MySecretKey", algo, plaintext)
	if err != nil {
		log.Fatalf("Second encryption failed: %v", err)
	}
	if bytes.Equal(ciphertext, ciphertext2) {
		log.Fatal("FATAL: Ciphertext should be different on each encryption!")
	}
	fmt.Println("   Second encryption successful (and produced different ciphertext due to random nonce).\n")

	// 6. Clean up
	fmt.Println("6. Cleaning up SDK resources...")
	Cleanup()

	fmt.Println("\n--- SDK Example Finished ---")
}
