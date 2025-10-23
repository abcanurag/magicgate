package com.example.sdk.internal;

import com.example.sdk.SDKException;

/**
 * Mock network client to simulate backend API calls.
 * In a real implementation, this would use java.net.http.HttpClient or another
 * HTTP client library to make real HTTPS requests.
 */
public class NetClient {

    private static final String MOCK_CONFIG_RESPONSE = "{\"api_version\":\"1.0\", \"features\":[\"AES-256-GCM\", \"RSA\"]}";
    private static final String MOCK_JWT_RESPONSE = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c";
    private static final String MOCK_KEY_RESPONSE = "0123456789abcdef0123456789abcdef"; // 32 bytes for AES-256

    public String fetchConfig(String regToken, String apiEndpoint) throws SDKException {
        System.out.printf("NET: Fetching config from '%s' with token '%s'%n", apiEndpoint, regToken);
        // Simulate network call
        return MOCK_CONFIG_RESPONSE;
    }

    public String authenticate(String identity, String secret, String apiEndpoint) throws SDKException {
        System.out.printf("NET: Authenticating user '%s' at '%s'%n", identity, apiEndpoint);
        // Simulate network call
        return MOCK_JWT_RESPONSE;
    }

    public String keyOp(String jwt, String opType, String keyName, byte[] keyData, String apiEndpoint) throws SDKException {
        System.out.printf("NET: Performing key op '%s' for key '%s' at '%s'%n", opType, keyName, apiEndpoint);
        // In a real implementation, the JWT would be in an "Authorization: Bearer <jwt>" header.

        switch (opType.toUpperCase()) {
            case "CREATE":
            case "UPDATE":
                System.out.printf("  -> Key Data: %s%n", keyData != null ? new String(keyData) : "N/A");
                return "{\"status\":\"success\"}";
            case "READ":
                return MOCK_KEY_RESPONSE;
            case "DELETE":
                return "{\"status\":\"success\"}";
            default:
                throw new SDKException("Unsupported key operation: " + opType);
        }
    }
}
