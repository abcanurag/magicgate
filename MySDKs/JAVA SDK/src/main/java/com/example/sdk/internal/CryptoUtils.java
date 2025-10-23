package com.example.sdk.internal;

import com.example.sdk.SDKException;

import javax.crypto.Cipher;
import javax.crypto.SecretKey;
import javax.crypto.spec.GCMParameterSpec;
import java.security.SecureRandom;

/**
 * Utility class for cryptographic operations using JCE.
 */
public final class CryptoUtils {

    private static final int GCM_IV_LENGTH = 12; // 96 bits
    private static final int GCM_TAG_LENGTH = 16; // 128 bits

    private CryptoUtils() {}

    public static byte[] encrypt(SecretKey key, String algoName, byte[] plaintext) throws SDKException {
        try {
            Cipher cipher = Cipher.getInstance(algoName);

            // Generate a random, non-repeating IV for each encryption
            byte[] iv = new byte[GCM_IV_LENGTH];
            new SecureRandom().nextBytes(iv);

            GCMParameterSpec gcmSpec = new GCMParameterSpec(GCM_TAG_LENGTH * 8, iv);
            cipher.init(Cipher.ENCRYPT_MODE, key, gcmSpec);

            byte[] ciphertext = cipher.doFinal(plaintext);

            // Prepend the IV to the ciphertext for use during decryption
            byte[] encryptedData = new byte[iv.length + ciphertext.length];
            System.arraycopy(iv, 0, encryptedData, 0, iv.length);
            System.arraycopy(ciphertext, 0, encryptedData, iv.length, ciphertext.length);

            return encryptedData;

        } catch (Exception e) {
            throw new SDKException("Encryption failed", e);
        }
    }

    public static byte[] decrypt(SecretKey key, String algoName, byte[] encryptedData) throws SDKException {
        try {
            Cipher cipher = Cipher.getInstance(algoName);

            // Extract the IV from the beginning of the encrypted data
            GCMParameterSpec gcmSpec = new GCMParameterSpec(GCM_TAG_LENGTH * 8, encryptedData, 0, GCM_IV_LENGTH);
            cipher.init(Cipher.DECRYPT_MODE, key, gcmSpec);

            // Decrypt the actual ciphertext (which comes after the IV)
            return cipher.doFinal(encryptedData, GCM_IV_LENGTH, encryptedData.length - GCM_IV_LENGTH);

        } catch (Exception e) {
            throw new SDKException("Decryption failed. The data may be corrupt or the key incorrect.", e);
        }
    }
}
