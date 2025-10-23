/*
### 2. Public Header: ``

This is the only header file that a third-party application needs to include. It defines the public API, error codes, and opaque structures.

```c
*/
#ifndef SDK_H
#define SDK_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief SDK Error Codes
 *
 * Defines a list of possible return codes for SDK functions.
 */
typedef enum {
    SDK_SUCCESS = 0,
    SDK_ERROR_GENERAL = -1,
    SDK_ERROR_NOT_INITIALIZED = -2,
    SDK_ERROR_ALREADY_INITIALIZED = -3,
    SDK_ERROR_INVALID_ARGUMENT = -4,
    SDK_ERROR_NETWORK = -5,
    SDK_ERROR_BACKEND_API = -6,
    SDK_ERROR_CRYPTO = -7,
    SDK_ERROR_BUFFER_TOO_SMALL = -8,
    SDK_ERROR_NO_SESSION = -9,
    SDK_ERROR_KEY_NOT_FOUND = -10,
    SDK_ERROR_UNSUPPORTED_OPERATION = -11,
    SDK_ERROR_MUTEX = -12
} SDK_STATUS;

/**
 * @brief Initializes the SDK.
 *
 * This function must be called once before any other SDK function. It registers the
 * SDK instance with the backend, fetches configuration, and initializes internal state.
 *
 * @param regToken A registration token for the application.
 * @return SDK_SUCCESS on success, or an SDK_STATUS error code on failure.
 */
int SDK_Init(const char *regToken);

/**
 * @brief Cleans up all resources used by the SDK.
 *
 * This function should be called when the application is shutting down to release
 * memory, close network connections, and destroy mutexes.
 */
void SDK_Cleanup(void);

/**
 * @brief Creates an authenticated session with the backend service.
 *
 * Authenticates using an identity and secret, receiving a JWT for subsequent
 * authenticated requests.
 *
 * @param identity The user or application identity.
 * @param secret The corresponding secret for the identity.
 * @param jwtBuffer Optional buffer to receive the JWT. Can be NULL.
 * @param jwtBufferLen The size of jwtBuffer.
 * @return SDK_SUCCESS on success, or an SDK_STATUS error code on failure.
 */
int SDK_CreateSession(const char *identity, const char *secret, char *jwtBuffer, size_t jwtBufferLen);

/**
 * @brief Performs a key management operation (CRUD).
 *
 * Interacts with the backend to create, read, update, or delete a key.
 *
 * @param opType The operation to perform: "CREATE", "READ", "UPDATE", "DELETE".
 * @param keyName The unique name of the key.
 * @param keyData For "CREATE" or "UPDATE", this is the key material. Can be NULL for other ops.
 * @return SDK_SUCCESS on success, or an SDK_STATUS error code on failure.
 */
int SDK_KeyOperation(const char *opType, const char *keyName, const char *keyData);

/**
 * @brief Performs a cryptographic operation using a managed key.
 *
 * Encrypts, decrypts, or signs data using a key specified by name. The key will be
 * fetched from the backend if not already cached locally.
 *
 * @param keyName The name of the key to use for the operation.
 * @param algoName The algorithm to use (e.g., "AES-256-GCM").
 * @param input The input data to be processed.
 * @param inputLen The length of the input data.
 * @param output Buffer to store the result (e.g., ciphertext, signature).
 * @param outputLen On input, the size of the output buffer. On success, the actual number of bytes written.
 * @return SDK_SUCCESS on success, or an SDK_STATUS error code on failure.
 */
int SDK_DoCrypto(const char *keyName, const char *algoName,
                   const unsigned char *input, size_t inputLen,
                   unsigned char *output, size_t *outputLen);

#ifdef __cplusplus
}
#endif

#endif // SDK_H
