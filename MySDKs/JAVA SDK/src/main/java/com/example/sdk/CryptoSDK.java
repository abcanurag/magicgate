package com.example.sdk;

import com.example.sdk.internal.SDKContext;
import com.example.sdk.internal.CryptoUtils;

import javax.crypto.SecretKey;
import java.util.concurrent.locks.ReentrantLock;

/**
 * Main public interface for the Crypto SDK.
 * This class provides static methods to initialize the SDK, manage sessions,
 * handle keys, and perform cryptographic operations.
 */
public final class CryptoSDK {

    private static SDKContext context;
    private static final ReentrantLock lock = new ReentrantLock();

    // Private constructor to prevent instantiation
    private CryptoSDK() {}

    /**
     * Initializes the SDK. Must be called once before any other operation.
     *
     * @param regToken The registration token for the application.
     * @throws SDKException if initialization fails or if already initialized.
     */
    public static void init(String regToken) throws SDKException {
        lock.lock();
        try {
            if (context != null && context.isInitialized()) {
                throw new SDKException("SDK is already initialized.");
            }
            context = new SDKContext();
            context.initialize(regToken);
            System.out.println("SDK Initialized successfully.");
        } finally {
            lock.unlock();
        }
    }

    /**
     * Cleans up all resources used by the SDK.
     */
    public static void cleanup() {
        lock.lock();
        try {
            if (context != null) {
                context.cleanup();
                context = null;
                System.out.println("SDK Cleanup complete.");
            }
        } finally {
            lock.unlock();
        }
    }

    /**
     * Creates an authenticated session with the backend service.
     *
     * @param identity The user or application identity.
     * @param secret   The corresponding secret.
     * @return The received JWT for the session.
     * @throws SDKException if the SDK is not initialized or authentication fails.
     */
    public static String createSession(String identity, String secret) throws SDKException {
        checkInitialized();
        return context.createSession(identity, secret);
    }

    /**
     * Performs a key management operation (CRUD).
     *
     * @param opType  The operation: "CREATE", "READ", "UPDATE", "DELETE".
     * @param keyName The unique name of the key.
     * @param keyData For "CREATE" or "UPDATE", the raw key material. Can be null.
     * @throws SDKException if the operation fails.
     */
    public static void keyOperation(String opType, String keyName, byte[] keyData) throws SDKException {
        checkInitialized();
        context.keyOperation(opType, keyName, keyData);
    }

    /**
     * Performs a cryptographic operation (encryption or decryption).
     *
     * @param opType   The crypto operation to perform ("ENCRYPT" or "DECRYPT").
     * @param keyName  The name of the key to use.
     * @param algoName The algorithm to use (e.g., "AES/GCM/NoPadding").
     * @param input    The data to be processed (plaintext or ciphertext).
     * @return The result of the operation (ciphertext or plaintext).
     * @throws SDKException if the operation fails.
     */
    public static byte[] doCrypto(String opType, String keyName, String algoName, byte[] input) throws SDKException {
        checkInitialized();

        // 1. Get key from context (fetches from backend if not in cache)
        SecretKey key = context.getKey(keyName);

        // 2. Perform crypto operation
        if ("ENCRYPT".equalsIgnoreCase(opType)) {
            return CryptoUtils.encrypt(key, algoName, input);
        } else if ("DECRYPT".equalsIgnoreCase(opType)) {
            return CryptoUtils.decrypt(key, algoName, input);
        } else {
            throw new SDKException("Unsupported crypto operation: " + opType);
        }
    }

    /**
     * Helper to ensure the SDK is initialized before use.
     */
    private static void checkInitialized() throws SDKException {
        if (context == null || !context.isInitialized()) {
            throw new SDKException("SDK is not initialized. Please call SDK.init() first.");
        }
    }
}
