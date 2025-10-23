package com.example.sdk.app;

import com.example.sdk.CryptoSDK;
import com.example.sdk.SDKException;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.HexFormat;

public class Example {

    public static void main(String[] args) {
        System.out.println("--- Java SDK Usage Example ---\n");

        try {
            // 1. Initialize the SDK
            System.out.println("1. Initializing SDK...");
            CryptoSDK.init("my-app-registration-token-12345");
            System.out.println();

            // 2. Create a session
            System.out.println("2. Creating a session...");
            String jwt = CryptoSDK.createSession("app_user_01", "super_secret_password");
            System.out.printf("   Session created. Received JWT: %.30s...%n%n", jwt);

            // 3. Manage a key (CREATE)
            System.out.println("3. Creating a key named 'MySecretKey'...");
            byte[] keyMaterial = "this-is-my-super-secret-key-data".getBytes(StandardCharsets.UTF_8);
            CryptoSDK.keyOperation("CREATE", "MySecretKey", keyMaterial);
            System.out.println("   Key 'MySecretKey' created on backend.\n");

            // 4. Perform crypto operations (Encrypt and Decrypt)
            System.out.println("4. Encrypting and decrypting data with 'MySecretKey'...");
            String plaintext = "This is a very sensitive message.";
            String algo = "AES/GCM/NoPadding";

            // Encrypt
            byte[] ciphertext = CryptoSDK.doCrypto("ENCRYPT", "MySecretKey", algo, plaintext.getBytes(StandardCharsets.UTF_8));
            System.out.println("   Encryption successful.");
            System.out.println("   Plaintext: '" + plaintext + "'");
            System.out.println("   Ciphertext (hex): " + HexFormat.of().formatHex(ciphertext) + "\n");

            // Decrypt
            byte[] decryptedTextBytes = CryptoSDK.doCrypto("DECRYPT", "MySecretKey", algo, ciphertext);
            String decryptedText = new String(decryptedTextBytes, StandardCharsets.UTF_8);
            System.out.println("   Decryption successful.");
            System.out.println("   Decrypted Text: '" + decryptedText + "'\n");

            // Verify
            assert plaintext.equals(decryptedText);
            System.out.println("   SUCCESS: Original plaintext matches decrypted text.\n");

            // 5. Demonstrate key caching
            System.out.println("5. Encrypting again (should use cached key)...");
            byte[] ciphertext2 = CryptoSDK.doCrypto("ENCRYPT", "MySecretKey", algo, plaintext.getBytes(StandardCharsets.UTF_8));
            assert !Arrays.equals(ciphertext, ciphertext2); // Should be different due to random IV
            System.out.println("   Second encryption successful.\n");

        } catch (SDKException e) {
            System.err.println("SDK operation failed: " + e.getMessage());
            e.printStackTrace();
        } finally {
            // 6. Clean up
            System.out.println("6. Cleaning up SDK resources...");
            CryptoSDK.cleanup();
        }

        System.out.println("\n--- SDK Example Finished ---");
    }
}
