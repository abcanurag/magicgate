package main

import (
	"fmt"
	"strings"
	"time"
)

// netClient simulates a network client.
type netClient struct {
	// In a real implementation, this would hold an http.Client
}

func newNetClient() *netClient {
	return &netClient{}
}

const (
	mockConfigResponse = `{"api_version":"1.0", "features":["AES-256-GCM"]}`
	mockJWTResponse    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	// 32 bytes for AES-256, represented as a hex string
	mockKeyResponse = "3031323334353637383961626364656630313233343536373839616263646566"
)

func (nc *netClient) FetchConfig(regToken, apiEndpoint string) (string, error) {
	fmt.Printf("NET: Fetching config from '%s' with token '%s'\n", apiEndpoint, regToken)
	time.Sleep(50 * time.Millisecond) // Simulate network latency
	return mockConfigResponse, nil
}

func (nc *netClient) Authenticate(identity, secret, apiEndpoint string) (string, error) {
	fmt.Printf("NET: Authenticating user '%s' at '%s'\n", identity, apiEndpoint)
	time.Sleep(50 * time.Millisecond)
	return mockJWTResponse, nil
}

func (nc *netClient) KeyOp(jwt, opType, keyName string, keyData []byte, apiEndpoint string) (string, error) {
	fmt.Printf("NET: Performing key op '%s' for key '%s' at '%s'\n", opType, keyName, apiEndpoint)
	time.Sleep(50 * time.Millisecond)

	switch strings.ToUpper(opType) {
	case "CREATE", "UPDATE":
		// In a real implementation, the JWT would be in an "Authorization: Bearer <jwt>" header.
		return `{"status":"success"}`, nil
	case "READ":
		return mockKeyResponse, nil
	case "DELETE":
		return `{"status":"success"}`, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedOperation, opType)
	}
}
