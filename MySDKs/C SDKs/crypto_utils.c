#include "sdk_internal.h"
#include <string.h>
#include <openssl/err.h>

/**
 * @brief Maps a string algorithm name to an OpenSSL EVP_CIPHER object.
 */
const EVP_CIPHER* crypto_get_evp_cipher(const char *algoName) {
    if (strcasecmp(algoName, "AES-256-GCM") == 0) {
        return EVP_aes_256_gcm();
    }
    // Add other mappings here, e.g., AES-256-CBC
    // if (strcasecmp(algoName, "AES-256-CBC") == 0) {
    //     return EVP_aes_256_cbc();
    // }
    return NULL;
}

/**
 * @brief A wrapper around OpenSSL's EVP encryption routines.
 * This example focuses on AES-GCM. A production implementation would handle
 * other modes and operations (decryption, signing, etc.).
 */
int crypto_encrypt(const EVP_CIPHER *cipher, const unsigned char *key, const unsigned char *iv,
                   const unsigned char *plaintext, int plaintext_len,
                   unsigned char *ciphertext, int *ciphertext_len) {
    EVP_CIPHER_CTX *ctx = NULL;
    int len;
    int final_len;
    int ret = SDK_ERROR_CRYPTO;

    // Create and initialize the context
    if (!(ctx = EVP_CIPHER_CTX_new())) goto cleanup;

    // Initialize encryption operation.
    if (1 != EVP_EncryptInit_ex(ctx, cipher, NULL, NULL, NULL)) goto cleanup;

    // Set IV length for GCM
    if (EVP_CIPHER_is_a(cipher, "GCM")) {
        if (1 != EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_SET_IVLEN, 12, NULL)) goto cleanup;
    }

    // Initialize key and IV
    if (1 != EVP_EncryptInit_ex(ctx, NULL, NULL, key, iv)) goto cleanup;

    // Provide the message to be encrypted, and obtain the encrypted output.
    if (1 != EVP_EncryptUpdate(ctx, ciphertext, &len, plaintext, plaintext_len)) goto cleanup;
    *ciphertext_len = len;

    // Finalize the encryption.
    if (1 != EVP_EncryptFinal_ex(ctx, ciphertext + len, &final_len)) goto cleanup;
    *ciphertext_len += final_len;

    // For GCM, get the tag
    if (EVP_CIPHER_is_a(cipher, "GCM")) {
        unsigned char tag[16];
        if (1 != EVP_CIPHER_CTX_ctrl(ctx, EVP_CTRL_GCM_GET_TAG, 16, tag)) goto cleanup;
        // Append the tag to the ciphertext. The recipient will need to know to
        // separate it for decryption/verification.
        memcpy(ciphertext + *ciphertext_len, tag, 16);
        *ciphertext_len += 16;
    }

    ret = SDK_SUCCESS;

cleanup:
    if (ret != SDK_SUCCESS) {
        ERR_print_errors_fp(stderr);
    }
    if (ctx) {
        EVP_CIPHER_CTX_free(ctx);
    }
    return ret;
}
