namespace MyCompany.CryptoSDK.Internal;

/// <summary>
/// Mock network client to simulate backend API calls.
/// In a real implementation, this would use HttpClient to make real HTTPS requests.
/// </summary>
internal class NetClient : IDisposable
{
    private static readonly string MOCK_CONFIG_RESPONSE = "{\"api_version\":\"1.0\", \"features\":[\"AES-256-GCM\"]}";
    private static readonly string MOCK_JWT_RESPONSE = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c";
    // 32 bytes for AES-256, represented as a hex string
    private static readonly string MOCK_KEY_RESPONSE = "3031323334353637383961626364656630313233343536373839616263646566"; 

    public async Task<string> FetchConfigAsync(string regToken, string apiEndpoint)
    {
        Console.WriteLine($"NET: Fetching config from '{apiEndpoint}' with token '{regToken}'");
        await Task.Delay(50); // Simulate network latency
        return MOCK_CONFIG_RESPONSE;
    }

    public async Task<string> AuthenticateAsync(string identity, string secret, string apiEndpoint)
    {
        Console.WriteLine($"NET: Authenticating user '{identity}' at '{apiEndpoint}'");
        await Task.Delay(50); // Simulate network latency
        return MOCK_JWT_RESPONSE;
    }

    public async Task<string> KeyOpAsync(string jwt, string opType, string keyName, byte[]? keyData, string apiEndpoint)
    {
        Console.WriteLine($"NET: Performing key op '{opType}' for key '{keyName}' at '{apiEndpoint}'");
        await Task.Delay(50); // Simulate network latency

        switch (opType.ToUpperInvariant())
        {
            case "CREATE":
            case "UPDATE":
                Console.WriteLine($"  -> Key Data (hex): { (keyData != null ? Convert.ToHexString(keyData) : "N/A") }");
                return "{\"status\":\"success\"}";
            case "READ":
                return MOCK_KEY_RESPONSE;
            case "DELETE":
                return "{\"status\":\"success\"}";
            default:
                throw new SDKException($"Unsupported key operation: {opType}");
        }
    }

    public void Dispose()
    {
        // In a real implementation with HttpClient, you might dispose it here if needed.
    }
}
