#ifndef SDK_INTERNAL_H
#define SDK_INTERNAL_H

#include "sdk.h"
#include <pthread.h>
#include <openssl/evp.h>

#define MAX_CONFIG_LEN 4096
#define MAX_JWT_LEN 1024
#define MAX_KEY_CACHE_SIZE 10
#define MAX_KEY_NAME_LEN 128
#define MAX_KEY_DATA_LEN 256

// Forward declarations for internal helper functions
int net_init(void);
void net_cleanup(void);
int net_fetch_config(const char *regToken, const char *api_endpoint, char *out_buffer, size_t buffer_len);
int net_authenticate(const char *identity, const char *secret, const char *api_endpoint, char *out_jwt, size_t jwt_len);
int net_key_op(const char *jwt, const char *opType, const char *keyName, const char *keyData, const char *api_endpoint, char *response_buf, size_t response_len);

const EVP_CIPHER* crypto_get_evp_cipher(const char *algoName);
int crypto_encrypt(const EVP_CIPHER *cipher, const unsigned char *key, const unsigned char *iv,
                   const unsigned char *plaintext, int plaintext_len,
                   unsigned char *ciphertext, int *ciphertext_len);

// Internal cache for cryptographic keys
typedef struct {
    char name[MAX_KEY_NAME_LEN];
    unsigned char data[MAX_KEY_DATA_LEN];
    size_t len;
} KeyCacheEntry;

/**
 * @brief The internal state of the SDK.
 *
 * This structure holds all configuration, session data, and handles needed
 * for the SDK to operate. It is managed as a singleton.
 */
typedef struct {
    int is_initialized;
    char config_json[MAX_CONFIG_LEN];
    char api_endpoint[256];
    char jwt[MAX_JWT_LEN];
    KeyCacheEntry key_cache[MAX_KEY_CACHE_SIZE];
    int key_cache_count;
    pthread_mutex_t lock;
} SDKContext;

#endif // SDK_INTERNAL_H
