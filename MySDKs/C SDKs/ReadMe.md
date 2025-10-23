You are an expert systems programmer and C SDK architect. 
Write a production-grade C SDK wrapper that extends OpenSSL and provides an abstraction layer to a backend crypto service. 

The SDK should be modular, easy to integrate into third-party applications, and expose the following API functions (exported symbols):

1. int SDK_Init(const char *regToken);
   - Registers the SDK instance with a backend server using a registration token.
   - Fetches configuration from the server corresponding to the application bound to the regToken.
   - Loads or caches the configuration locally (JSON or key-value format).
   - Initializes OpenSSL contexts and any internal state.

2. int SDK_CreateSession(const char *identity, const char *secret, char *jwtBuffer, size_t jwtBufferLen);
   - Creates an authenticated session by sending identity credentials to the server.
   - Receives and caches a JWT token internally (optionally returns it to caller in jwtBuffer).
   - JWT will be used for subsequent authenticated requests (CRUD or crypto operations).

3. int SDK_KeyOperation(const char *opType, const char *keyName, const char *keyData);
   - Implements CRUD-like operations for key management:
     opType can be "CREATE", "READ", "UPDATE", or "DELETE".
   - Interacts with the backend to manage keys associated with the current session/app.
   - May store keys locally (secure enclave or memory buffer).

4. int SDK_DoCrypto(const char *keyName, const char *algoName, 
                    const unsigned char *input, size_t inputLen,
                    unsigned char *output, size_t *outputLen);
   - Performs encryption/decryption or signing operations using the specified key.
   - Internally uses OpenSSL EVP APIs (e.g., EVP_EncryptInit_ex, EVP_EncryptUpdate, EVP_EncryptFinal_ex).
   - Select the algorithm dynamically based on algoName (e.g., AES-256-GCM, RSA, ECDSA).
   - Returns ciphertext or signature in `output`.

Additional Requirements:
- All functions should return well-defined error codes (enum-based).
- Implement proper resource cleanup with an SDK_Cleanup() function.
- Include an internal state/context struct (e.g., SDKContext) to hold session data, config, and OpenSSL handles.
- Demonstrate thread safety using mutexes where necessary.
- Provide minimal stub/mock backend API calls (e.g., via HTTPS using libcurl or similar).
- The code should compile cleanly with gcc/clang on Linux and be extensible for Windows.
- Include brief usage examples showing how a caller application would initialize, create a session, perform key ops, and encrypt data.

Focus on writing modular, well-commented code with clear separation between:
- SDK Public Header (`sdk.h`)
- SDK Implementation (`sdk.c`)
- Internal Helper Modules (`net_client.c`, `crypto_utils.c`)

Use OpenSSL’s EVP interface for crypto and provide clear error handling.

Generate all necessary C files, headers, and a short README snippet describing build instructions (e.g., `gcc -o testsdk sdk.c net_client.c crypto_utils.c -lssl -lcrypto -lcurl`).
