package com.example.sdk;

/**
 * Custom exception for all errors originating from the SDK.
 */
public class SDKException extends Exception {
    public SDKException(String message) {
        super(message);
    }

    public SDKException(String message, Throwable cause) {
        super(message, cause);
    }
}
