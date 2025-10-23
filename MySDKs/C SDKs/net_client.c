#include "sdk_internal.h"
#include <stdio.h>
#include <string.h>
#include <curl/curl.h>

// Mock backend responses
static const char* MOCK_CONFIG_RESPONSE = "{\"api_version\":\"1.0\", \"features\":[\"AES-256-GCM\", \"RSA\"]}";
static const char* MOCK_JWT_RESPONSE = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c";
static const char* MOCK_KEY_RESPONSE = "0123456789abcdef0123456789abcdef"; // 32 bytes for AES-256

int net_init(void) {
    return (curl_global_init(CURL_GLOBAL_DEFAULT) == 0) ? 0 : -1;
}

void net_cleanup(void) {
    curl_global_cleanup();
}

/**
 * @brief Mock function to fetch configuration.
 * In a real implementation, this would make an HTTPS GET request.
 */
int net_fetch_config(const char *regToken, const char *api_endpoint, char *out_buffer, size_t buffer_len) {
    printf("NET: Fetching config from '%s' with token '%s'\n", api_endpoint, regToken);
    
    if (strlen(MOCK_CONFIG_RESPONSE) + 1 > buffer_len) {
        return SDK_ERROR_BUFFER_TOO_SMALL;
    }
    strcpy(out_buffer, MOCK_CONFIG_RESPONSE);
    
    // Simulate network latency
    // sleep(1);
    
    return SDK_SUCCESS;
}

/**
 * @brief Mock function to authenticate and get a JWT.
 * In a real implementation, this would make an HTTPS POST request.
 */
int net_authenticate(const char *identity, const char *secret, const char *api_endpoint, char *out_jwt, size_t jwt_len) {
    printf("NET: Authenticating user '%s' at '%s'\n", identity, api_endpoint);
    (void)secret; // Unused in mock

    if (strlen(MOCK_JWT_RESPONSE) + 1 > jwt_len) {
        return SDK_ERROR_BUFFER_TOO_SMALL;
    }
    strcpy(out_jwt, MOCK_JWT_RESPONSE);
    return SDK_SUCCESS;
}

/**
 * @brief Mock function for key operations.
 * In a real implementation, this would make authenticated HTTPS requests (POST, GET, DELETE).
 */
int net_key_op(const char *jwt, const char *opType, const char *keyName, const char *keyData, const char *api_endpoint, char *response_buf, size_t response_len) {
    printf("NET: Performing key op '%s' for key '%s' at '%s'\n", opType, keyName, api_endpoint);
    (void)jwt; // Unused in mock, but would be in Authorization header

    if (strcasecmp(opType, "CREATE") == 0) {
        printf("  -> Key Data: %s\n", keyData);
        // Backend would store this key
    } else if (strcasecmp(opType, "READ") == 0) {
        if (strlen(MOCK_KEY_RESPONSE) + 1 > response_len) {
            return SDK_ERROR_BUFFER_TOO_SMALL;
        }
        strcpy(response_buf, MOCK_KEY_RESPONSE);
    } else if (strcasecmp(opType, "UPDATE") == 0) {
        printf("  -> New Key Data: %s\n", keyData);
    } else if (strcasecmp(opType, "DELETE") == 0) {
        // Backend would delete this key
    } else {
        return SDK_ERROR_UNSUPPORTED_OPERATION;
    }
    
    return SDK_SUCCESS;
}
