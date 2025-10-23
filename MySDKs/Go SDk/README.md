# Backend Crypto Go SDK

This SDK provides a Go-language wrapper for interacting with a backend crypto service. It abstracts away the complexities of direct API calls, session management, and cryptographic operations, offering a simple, idiomatic, and secure interface for client applications.

The SDK uses Go's standard `crypto` library for cryptographic primitives and the built-in `net/http` package for network communication.

## Features

-   **Idiomatic Go API**: Simple functions and error handling.
-   **Initialization**: Register the SDK instance and fetch configuration from the backend.
-   **Session Management**: Create authenticated sessions using JWTs.
-   **Key Management**: Perform CRUD (Create, Read, Update, Delete) operations on cryptographic keys.
-   **Cryptographic Operations**: Perform high-level encryption/decryption using keys managed by the service.
-   **Thread-Safe**: Designed for concurrent use in goroutines.

## Building and Running the Example

### Prerequisites

You will need the Go toolchain (version 1.18 or newer) installed on your system.

### Run the Example

Navigate to the root directory of the project (where the `.go` files are located) and run the following command. The `go` tool will automatically handle dependencies, compilation, and execution.

```sh
go run .
