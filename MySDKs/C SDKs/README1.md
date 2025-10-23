# Backend Crypto C SDK

This SDK provides a C-language wrapper for interacting with a backend crypto service. It abstracts away the complexities of direct API calls, session management, and cryptographic operations, offering a simple and secure interface for client applications.

The SDK uses OpenSSL for underlying cryptographic primitives and `libcurl` for network communication.

## Features

-   **Initialization**: Register the SDK instance and fetch configuration from the backend.
-   **Session Management**: Create authenticated sessions using JWTs.
-   **Key Management**: Perform CRUD (Create, Read, Update, Delete) operations on cryptographic keys.
-   **Cryptographic Operations**: Perform high-level encryption/decryption using keys managed by the service.
-   **Thread-Safe**: Designed for use in multi-threaded applications.

## Building the SDK and Example

You will need `libssl`, `libcrypto`, and `libcurl` development libraries installed on your system.

-   On Debian/Ubuntu: `sudo apt-get install libssl-dev libcurl4-openssl-dev`
-   On RHEL/CentOS: `sudo yum install openssl-devel libcurl-devel`

To compile the example usage file which links against the SDK source files, use the following command:

```sh
gcc -o example_app example.c sdk.c net_client.c crypto_utils.c -lssl -lcrypto -lcurl -lpthread -Wall -Wextra -g

