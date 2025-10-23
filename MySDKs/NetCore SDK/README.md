# Backend Crypto .NET SDK

This SDK provides a C#/.NET wrapper for interacting with a backend crypto service. It abstracts away the complexities of direct API calls, session management, and cryptographic operations, offering a simple, modern, and secure interface for client applications.

The SDK uses the standard .NET cryptography libraries (`System.Security.Cryptography`) and the built-in `HttpClient` for network communication.

## Features

-   **Asynchronous API**: All network-bound operations are `async` for better performance in modern applications.
-   **Initialization**: Register the SDK instance and fetch configuration from the backend.
-   **Session Management**: Create authenticated sessions using JWTs.
-   **Key Management**: Perform CRUD (Create, Read, Update, Delete) operations on cryptographic keys.
-   **Cryptographic Operations**: Perform high-level encryption/decryption using keys managed by the service.
-   **Thread-Safe**: Designed for use in multi-threaded applications.
-   **Exception-Based Error Handling**: Uses custom exceptions for clear error reporting.

## Building and Running the Example

This project is set up to be built with the `dotnet` CLI (.NET 6.0 or newer).

### Prerequisites

You will need the .NET SDK (version 6.0 or later) installed on your system.

### Compile and Run

1.  **Restore dependencies and build the project:**
    Navigate to the root directory of the project (where the `.csproj` file is located) and run:
    ```sh
    dotnet build
    ```

2.  **Run the example application:**
    After a successful build, run the example project:
    ```sh
    dotnet run
    ```
