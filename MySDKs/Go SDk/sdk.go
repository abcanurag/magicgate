package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// SDK-specific errors
var (
	ErrSDKAlreadyInitialized    = errors.New("SDK is already initialized")
	ErrSDKNotInitialized        = errors.New("SDK is not initialized")
	ErrNoActiveSession          = errors.New("no active session, please call CreateSession first")
	ErrKeyNotFound              = errors.New("key not found")
	ErrUnsupportedOperation     = errors.New("unsupported operation")
	ErrRegistrationTokenInvalid = errors.New("registration token cannot be empty")
)

// sdkContext holds the internal state of the SDK.
type sdkContext struct {
	isInitialized bool
	apiEndpoint   string
	jwt           string
	keyCache      map[string][]byte
	lock          sync.RWMutex
	netClient     *netClient
}

var (
	context *sdkContext
	once    sync.Once
)

// Init initializes the SDK. Must be called once before any other operation.
func Init(regToken string) error {
	var initErr error
	once.Do(func() {
		if regToken == "" {
			initErr = ErrRegistrationTokenInvalid
			return
		}

		context = &sdkContext{
			keyCache:  make(map[string][]byte),
			netClient: newNetClient(),
		}

		// Mock backend endpoint
		context.apiEndpoint = "https://api.example-crypto.com/v1"

		// Fetch and parse configuration
		configJSON, err := context.netClient.FetchConfig(regToken, context.apiEndpoint)
		if err != nil {
			initErr = fmt.Errorf("failed to fetch config: %w", err)
			return
		}

		var configMap map[string]any
		if err := json.Unmarshal([]byte(configJSON), &configMap); err != nil {
			initErr = fmt.Errorf("failed to parse config: %w", err)
			return
		}
		fmt.Printf("SDK_Init: Configuration fetched: %s\n", configJSON)

		context.isInitialized = true
	})

	if context == nil || !context.isInitialized {
		// Reset once if initialization failed, allowing another attempt.
		if initErr != nil {
			once = sync.Once{}
			context = nil
		}
		return initErr
	}

	// If Init is called a second time after successful initialization
	if initErr == nil && context.isInitialized {
		// This part of the check is tricky with once.Do. A simple check is better.
	}
	// A simpler approach might be a global lock instead of sync.Once for re-initialization logic.
	// For this example, we assume Init is called correctly once.

	return nil
}

// Cleanup cleans up all resources used by the SDK.
func Cleanup() {
	if context == nil || !context.isInitialized {
		return
	}
	context.lock.Lock()
	defer context.lock.Unlock()

	// Securely clear sensitive data
	context.jwt = ""
	for k, v := range context.keyCache {
		for i := range v {
			v[i] = 0
		}
		delete(context.keyCache, k)
	}
	context.isInitialized = false
	fmt.Println("SDK Cleanup complete.")
}

// CreateSession creates an authenticated session with the backend service.
func CreateSession(identity, secret string) (string, error) {
	if context == nil || !context.isInitialized {
		return "", ErrSDKNotInitialized
	}
	context.lock.Lock()
	defer context.lock.Unlock()

	jwt, err := context.netClient.Authenticate(identity, secret, context.apiEndpoint)
	if err != nil {
		return "", err
	}
	context.jwt = jwt
	return jwt, nil
}

// KeyOperation performs a key management operation (CRUD).
func KeyOperation(opType, keyName string, keyData []byte) error {
	if context == nil || !context.isInitialized {
		return ErrSDKNotInitialized
	}

	context.lock.RLock()
	jwt := context.jwt
	context.lock.RUnlock()

	if jwt == "" {
		return ErrNoActiveSession
	}

	response, err := context.netClient.KeyOp(jwt, opType, keyName, keyData, context.apiEndpoint)
	if err != nil {
		return err
	}

	op := strings.ToUpper(opType)
	if op == "READ" {
		// Assuming the response is the raw key material in hex
		rawKey, err := hex.DecodeString(response)
		if err != nil {
			return fmt.Errorf("failed to decode key from hex: %w", err)
		}
		context.lock.Lock()
		context.keyCache[keyName] = rawKey
		context.lock.Unlock()
	} else if op == "DELETE" {
		context.lock.Lock()
		delete(context.keyCache, keyName)
		context.lock.Unlock()
	}

	return nil
}

// DoCrypto performs a cryptographic operation (encryption or decryption).
func DoCrypto(opType, keyName, algoName string, input []byte) ([]byte, error) {
	if context == nil || !context.isInitialized {
		return nil, ErrSDKNotInitialized
	}

	// 1. Get key from context (fetches from backend if not in cache)
	key, err := getKey(keyName)
	if err != nil {
		return nil, err
	}

	// 2. Perform crypto operation
	op := strings.ToUpper(opType)
	if op == "ENCRYPT" {
		return encrypt(key, algoName, input)
	}
	if op == "DECRYPT" {
		return decrypt(key, algoName, input)
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedOperation, opType)
}

// getKey is an internal helper to retrieve a key, fetching if not cached.
func getKey(keyName string) ([]byte, error) {
	context.lock.RLock()
	key, ok := context.keyCache[keyName]
	context.lock.RUnlock()

	if ok {
		return key, nil
	}

	// If not in cache, fetch from backend
	fmt.Printf("Key '%s' not in cache. Fetching from backend...\n", keyName)
	if err := KeyOperation("READ", keyName, nil); err != nil {
		return nil, err
	}

	// Try getting from cache again
	context.lock.RLock()
	key, ok = context.keyCache[keyName]
	context.lock.RUnlock()

	if !ok {
		return nil, ErrKeyNotFound
	}
	return key, nil
}
