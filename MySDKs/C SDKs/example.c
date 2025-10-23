#include "sdk.h"
#include <stdio.h>
#include <string.h>
#include <assert.h>

void print_hex(const char* label, const unsigned char* data, size_t len) {
    printf("%s (%zu bytes): ", label, len);
    for (size_t i = 0; i < len; ++i) {
        printf("%02x", data[i]);
    }
    printf("\n");
}

int main() {
    printf("--- SDK Usage Example ---\n\n");

    // 1. Initialize the SDK
    printf("1. Initializing SDK...\n");
    int status = SDK_Init("my-app-registration-token-12345");
    if (status != SDK_SUCCESS) {
        fprintf(stderr, "Failed to initialize SDK, error: %d\n", status);
        return 1;
    }
    printf("   SDK Initialized successfully.\n\n");

    // 2. Create a session
    printf("2. Creating a session...\n");
    char jwt_buffer[512];
    status = SDK_CreateSession("app_user_01", "super_secret_password", jwt_buffer, sizeof(jwt_buffer));
    if (status != SDK_SUCCESS) {
        fprintf(stderr, "Failed to create session, error: %d\n", status);
        SDK_Cleanup();
        return 1;
    }
    printf("   Session created. Received JWT: %.30s...\n\n", jwt_buffer);

    // 3. Manage a key (CREATE)
    printf("3. Creating a key named 'MySecretKey'...\n");
    // In a real scenario, the keyData might be generated locally and sent to the backend for storage.
    const char* key_material = "this-is-my-super-secret-key-data";
    status = SDK_KeyOperation("CREATE", "MySecretKey", key_material);
    if (status != SDK_SUCCESS) {
        fprintf(stderr, "Failed to create key, error: %d\n", status);
    } else {
        printf("   Key 'MySecretKey' created on backend.\n\n");
    }

    // 4. Perform a crypto operation
    printf("4. Encrypting data with 'MySecretKey'...\n");
    const char* plaintext = "This is a very sensitive message.";
    unsigned char ciphertext[256];
    size_t ciphertext_len = sizeof(ciphertext);

    status = SDK_DoCrypto("MySecretKey", "AES-256-GCM",
                          (const unsigned char*)plaintext, strlen(plaintext),
                          ciphertext, &ciphertext_len);

    if (status != SDK_SUCCESS) {
        fprintf(stderr, "Failed to perform crypto operation, error: %d\n", status);
    } else {
        printf("   Encryption successful.\n");
        printf("   Plaintext: '%s'\n", plaintext);
        print_hex("   Ciphertext", ciphertext, ciphertext_len);
        printf("\n");
    }
    
    // 5. Demonstrate key caching
    printf("5. Encrypting again (should use cached key)...\n");
    status = SDK_DoCrypto("MySecretKey", "AES-256-GCM",
                          (const unsigned char*)plaintext, strlen(plaintext),
                          ciphertext, &ciphertext_len);
    assert(status == SDK_SUCCESS);
    printf("   Second encryption successful.\n\n");


    // 6. Clean up
    printf("6. Cleaning up SDK resources...\n");
    SDK_Cleanup();
    printf("   Cleanup complete.\n\n");
    
    printf("--- SDK Example Finished ---\n");

    return 0;
}
