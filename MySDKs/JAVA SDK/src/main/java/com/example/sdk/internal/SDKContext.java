package com.example.sdk.internal;

import com.example.sdk.SDKException;
import com.google.gson.Gson;

import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantLock;

/**
 * Manages the internal state of the SDK, including configuration,
 * session data, and a thread-safe key cache.
 */
public class SDKContext {
    private volatile boolean initialized = false;
    private String apiEndpoint;
    private String jwt;
    private final Map<String, SecretKey> keyCache = new ConcurrentHashMap<>();
    private final ReentrantLock sessionLock = new ReentrantLock();
    private final NetClient netClient;

    public SDKContext() {
        this.netClient = new NetClient();
    }

    public void initialize(String regToken) throws SDKException {
        if (regToken == null || regToken.isEmpty()) {
            throw new SDKException("Registration token cannot be null or empty.");
        }

        // Mock backend endpoint
        this.apiEndpoint = "https://api.example-crypto.com/v1";

        // Fetch and parse configuration
        String configJson = netClient.fetchConfig(regToken, this.apiEndpoint);
        // Using Gson to parse config. In a real app, you'd map this to a config object.
        Gson gson = new Gson();
        Map<?, ?> configMap = gson.fromJson(configJson, Map.class);
        System.out.println("SDK_Init: Configuration fetched: " + configMap);

        this.initialized = true;
    }

    public String createSession(String identity, String secret) throws SDKException {
        sessionLock.lock();
        try {
            this.jwt = netClient.authenticate(identity, secret, this.apiEndpoint);
            return this.jwt;
        } finally {
            sessionLock.unlock();
        }
    }

    public void keyOperation(String opType, String keyName, byte[] keyData) throws SDKException {
        if (this.jwt == null) {
            throw new SDKException("No active session. Please call createSession() first.");
        }

        String response = netClient.keyOp(this.jwt, opType, keyName, keyData, this.apiEndpoint);

        if ("READ".equalsIgnoreCase(opType)) {
            // Assuming the response is the raw key material
            byte[] rawKey = response.getBytes(); // In a real scenario, this would be Base64 decoded
            SecretKey secretKey = new SecretKeySpec(rawKey, 0, rawKey.length, "AES");
            keyCache.put(keyName, secretKey);
        } else if ("DELETE".equalsIgnoreCase(opType)) {
            keyCache.remove(keyName);
        }
    }

    public SecretKey getKey(String keyName) throws SDKException {
        // Check cache first
        SecretKey key = keyCache.get(keyName);
        if (key != null) {
            return key;
        }

        // If not in cache, fetch from backend
        System.out.println("Key '" + keyName + "' not in cache. Fetching from backend...");
        keyOperation("READ", keyName, null);

        // Try getting from cache again
        key = keyCache.get(keyName);
        if (key == null) {
            throw new SDKException("Key not found: " + keyName);
        }
        return key;
    }

    public boolean isInitialized() {
        return initialized;
    }

    public void cleanup() {
        // Securely clear sensitive data
        this.jwt = null;
        keyCache.clear();
        this.initialized = false;
    }
}
