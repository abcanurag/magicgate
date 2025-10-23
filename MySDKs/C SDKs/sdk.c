#define _GNU_SOURCE // For strcasecmp
#include "sdk_internal.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <openssl/err.h>
#include <openssl/evp.h>

// Global SDK context singleton
static SDKContext g_ctx;

// --- Internal Helper Functions ---

/**
 * @brief Finds a key in the local cache.
 * @note This function must be called with the context lock held.
 */
static KeyCacheEntry* find_key_in_cache(const char *keyName) {
    for (int i = 0; i < g_ctx.key_cache_count; ++i) {
        if (strncmp(g_ctx.key_cache[i].name, keyName, MAX_KEY_NAME_LEN) == 0) {
            return &g_ctx.key_cache[i];
        }
    }
    return NULL;
}

/**
 * @brief Adds a key to the local cache.
 * @note This function must be called with the context lock held.
 */
static int add_key_to_cache(const char *keyName, const unsigned char *keyData, size_t keyLen) {
    if (g_ctx.key_cache_count >= MAX_KEY_CACHE_SIZE) {
        // Simple eviction: remove the oldest key
        memmove(&g_ctx.key_cache[0], &g_ctx.key_cache[1], sizeof(KeyCacheEntry) * (MAX_KEY_CACHE_SIZE - 1));
        g_ctx.key_cache_count--;
    }

    KeyCacheEntry *entry = &g_ctx.key_cache[g_ctx.key_cache_count];
    strncpy(entry->name, keyName, MAX_KEY_NAME_LEN - 1);
    entry->name[MAX_KEY_NAME_LEN - 1] = '\0';
    
    if (keyLen > MAX_KEY_DATA_LEN) return SDK_ERROR_BUFFER_TOO_SMALL;
    memcpy(entry->data, keyData, keyLen);
    entry->len = keyLen;
    
    g_ctx.key_cache_count++;
    return SDK_SUCCESS;
}

// --- Public API Implementation ---

int SDK_Init(const char *regToken) {
    if (g_ctx.is_initialized) {
        return SDK_ERROR_ALREADY_INITIALIZED;
    }
    if (!regToken || *regToken == '\0') {
        return SDK_ERROR_INVALID_ARGUMENT;
    }

    memset(&g_ctx, 0, sizeof(SDKContext));

    if (pthread_mutex_init(&g_ctx.lock, NULL) != 0) {
        return SDK_ERROR_MUTEX;
    }

    // Initialize OpenSSL
    OpenSSL_add_all_algorithms();
    ERR_load_crypto_strings();

    if (net_init() != 0) {
        pthread_mutex_destroy(&g_ctx.lock);
        return SDK_ERROR_NETWORK;
    }
    
    // Mock backend endpoint
    strncpy(g_ctx.api_endpoint, "https://api.example-crypto.com/v1", sizeof(g_ctx.api_endpoint) - 1);

    // Fetch configuration from backend
    if (net_fetch_config(regToken, g_ctx.api_endpoint, g_ctx.config_json, MAX_CONFIG_LEN) != 0) {
        net_cleanup();
        pthread_mutex_destroy(&g_ctx.lock);
        return SDK_ERROR_NETWORK;
    }

    // Here you would parse the JSON config, e.g., using a library like cJSON
    // For this example, we assume the config is simple and valid.
    printf("SDK_Init: Configuration fetched: %s\n", g_ctx.config_json);

    g_ctx.is_initialized = 1;
    return SDK_SUCCESS;
}

void SDK_Cleanup(void) {
    if (!g_ctx.is_initialized) {
        return;
    }
    pthread_mutex_lock(&g_ctx.lock);

    net_cleanup();
    EVP_cleanup();
    ERR_free_strings();
    
    // Securely zero out sensitive data
    memset(&g_ctx, 0, sizeof(SDKContext));
    
    g_ctx.is_initialized = 0; // Mark as uninitialized *after* zeroing
    
    pthread_mutex_unlock(&g_ctx.lock);
    pthread_mutex_destroy(&g_ctx.lock);
}

int SDK_CreateSession(const char *identity, const char *secret, char *jwtBuffer, size_t jwtBufferLen) {
    if (!g_ctx.is_initialized) return SDK_ERROR_NOT_INITIALIZED;
    if (!identity || !secret) return SDK_ERROR_INVALID_ARGUMENT;

    char temp_jwt[MAX_JWT_LEN];
    if (net_authenticate(identity, secret, g_ctx.api_endpoint, temp_jwt, MAX_JWT_LEN) != 0) {
        return SDK_ERROR_BACKEND_API;
    }

    pthread_mutex_lock(&g_ctx.lock);
    strncpy(g_ctx.jwt, temp_jwt, MAX_JWT_LEN - 1);
    g_ctx.jwt[MAX_JWT_LEN - 1] = '\0';
    pthread_mutex_unlock(&g_ctx.lock);

    if (jwtBuffer && jwtBufferLen > 0) {
        if (strlen(temp_jwt) + 1 > jwtBufferLen) {
            return SDK_ERROR_BUFFER_TOO_SMALL;
        }
        strcpy(jwtBuffer, temp_jwt);
    }

    return SDK_SUCCESS;
}

int SDK_KeyOperation(const char *opType, const char *keyName, const char *keyData) {
    if (!g_ctx.is_initialized) return SDK_ERROR_NOT_INITIALIZED;
    if (!opType || !keyName) return SDK_ERROR_INVALID_ARGUMENT;

    pthread_mutex_lock(&g_ctx.lock);
    if (g_ctx.jwt[0] == '\0') {
        pthread_mutex_unlock(&g_ctx.lock);
        return SDK_ERROR_NO_SESSION;
    }
    
    char current_jwt[MAX_JWT_LEN];
    strncpy(current_jwt, g_ctx.jwt, MAX_JWT_LEN);
    pthread_mutex_unlock(&g_ctx.lock);

    char response_buf[MAX_KEY_DATA_LEN] = {0};
    int status = net_key_op(current_jwt, opType, keyName, keyData, g_ctx.api_endpoint, response_buf, sizeof(response_buf));
    if (status != SDK_SUCCESS) {
        return status;
    }

    if (strcasecmp(opType, "READ") == 0) {
        pthread_mutex_lock(&g_ctx.lock);
        // Assuming response_buf contains the raw key bytes
        add_key_to_cache(keyName, (unsigned char*)response_buf, strlen(response_buf));
        pthread_mutex_unlock(&g_ctx.lock);
    } else if (strcasecmp(opType, "DELETE") == 0) {
        pthread_mutex_lock(&g_ctx.lock);
        KeyCacheEntry* entry = find_key_in_cache(keyName);
        if (entry) {
            // Invalidate cache entry. A more robust implementation would shift the array.
            memset(entry, 0, sizeof(KeyCacheEntry));
        }
        pthread_mutex_unlock(&g_ctx.lock);
    }

    return SDK_SUCCESS;
}

int SDK_DoCrypto(const char *keyName, const char *algoName,
                   const unsigned char *input, size_t inputLen,
                   unsigned char *output, size_t *outputLen) {
    if (!g_ctx.is_initialized) return SDK_ERROR_NOT_INITIALIZED;
    if (!keyName || !algoName || !input || !output || !outputLen || *outputLen == 0) {
        return SDK_ERROR_INVALID_ARGUMENT;
    }

    pthread_mutex_lock(&g_ctx.lock);
    if (g_ctx.jwt[0] == '\0') {
        pthread_mutex_unlock(&g_ctx.lock);
        return SDK_ERROR_NO_SESSION;
    }
    pthread_mutex_unlock(&g_ctx.lock);

    // 1. Select algorithm
    const EVP_CIPHER *cipher = crypto_get_evp_cipher(algoName);
    if (!cipher) {
        return SDK_ERROR_UNSUPPORTED_OPERATION;
    }

    // 2. Find key in cache
    pthread_mutex_lock(&g_ctx.lock);
    KeyCacheEntry* key_entry = find_key_in_cache(keyName);
    pthread_mutex_unlock(&g_ctx.lock);

    // 3. If not in cache, fetch from backend
    if (!key_entry) {
        printf("Key '%s' not in cache. Fetching from backend...\n", keyName);
        if (SDK_KeyOperation("READ", keyName, NULL) != SDK_SUCCESS) {
            return SDK_ERROR_KEY_NOT_FOUND;
        }
        // Try finding it again
        pthread_mutex_lock(&g_ctx.lock);
        key_entry = find_key_in_cache(keyName);
        pthread_mutex_unlock(&g_ctx.lock);
        if (!key_entry) {
            return SDK_ERROR_KEY_NOT_FOUND; // Should not happen if READ succeeded
        }
    }

    // 4. Perform crypto operation
    // For GCM, we need an IV. A production system would generate a unique IV per encryption.
    // Here we use a fixed IV for simplicity.
    unsigned char iv[12] = "my-unique-iv";
    int ciphertext_len;
    
    // We need a copy of the key data because key_entry is in shared memory
    unsigned char key_copy[MAX_KEY_DATA_LEN];
    size_t key_len_copy;
    pthread_mutex_lock(&g_ctx.lock);
    memcpy(key_copy, key_entry->data, key_entry->len);
    key_len_copy = key_entry->len;
    pthread_mutex_unlock(&g_ctx.lock);

    // Check key length against cipher's requirement
    if (EVP_CIPHER_key_length(cipher) != key_len_copy) {
        fprintf(stderr, "Error: Invalid key length for %s. Expected %d, got %zu.\n",
                algoName, EVP_CIPHER_key_length(cipher), key_len_copy);
        return SDK_ERROR_CRYPTO;
    }

    int status = crypto_encrypt(cipher, key_copy, iv, input, inputLen, output, &ciphertext_len);
    
    if (status != SDK_SUCCESS) {
        return status;
    }

    if ((size_t)ciphertext_len > *outputLen) {
        *outputLen = ciphertext_len;
        return SDK_ERROR_BUFFER_TOO_SMALL;
    }

    *outputLen = ciphertext_len;
    return SDK_SUCCESS;
}
