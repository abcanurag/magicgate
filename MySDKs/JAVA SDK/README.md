# Backend Crypto Java SDK

This SDK provides a Java-language wrapper for interacting with a backend crypto service. It abstracts away the complexities of direct API calls, session management, and cryptographic operations, offering a simple and secure interface for client applications.

The SDK uses the standard Java Cryptography Extension (JCE) for cryptographic primitives and a standard HTTP client for network communication.

## Features

-   **Initialization**: Register the SDK instance and fetch configuration from the backend.
-   **Session Management**: Create authenticated sessions using JWTs.
-   **Key Management**: Perform CRUD (Create, Read, Update, Delete) operations on cryptographic keys.
-   **Cryptographic Operations**: Perform high-level encryption/decryption using keys managed by the service.
-   **Thread-Safe**: Designed for use in multi-threaded applications.
-   **Exception-Based Error Handling**: Uses custom exceptions for clear error reporting.

## Building and Running the Example

This project is set up to be built with [Apache Maven](https://maven.apache.org/).

### Dependencies

The project uses the following external libraries, which Maven will download automatically:

-   `com.google.code.gson:gson`: For parsing JSON configuration from the backend.

### Compile and Run

1.  **Compile the project:**
    Navigate to the root directory of the project (where `pom.xml` is located) and run:
    ```sh
    mvn compile
    ```

2.  **Run the example application:**
    After a successful compilation, run the example class:
    ```sh
    mvn exec:java -Dexec.mainClass="com.example.sdk.app.Example"
    ```
