using System.Collections.Concurrent;
using System.Text.Json;

namespace MyCompany.CryptoSDK.Internal;

/// <summary>
/// Manages the internal state of the SDK, including configuration,
/// session data, and a thread-safe key cache.
/// </summary>
internal class SDKContext : IDisposable
{
    private string? _apiEndpoint;
    private string? _jwt;
    private readonly ConcurrentDictionary<string, byte[]> _keyCache = new();
    private readonly NetClient _netClient;

    public bool IsInitialized { get; private set; }

    public SDKContext()
    {
        _netClient = new NetClient();
    }

    public async Task InitializeAsync(string regToken)
    {
        if (string.IsNullOrEmpty(regToken))
        {
            throw new SDKException("Registration token cannot be null or empty.");
        }

        // Mock backend endpoint
        _apiEndpoint = "https://api.example-crypto.com/v1";

        // Fetch and parse configuration
        string configJson = await _netClient.FetchConfigAsync(regToken, _apiEndpoint);
        
        // Using System.Text.Json to parse config.
        var configMap = JsonSerializer.Deserialize<Dictionary<string, object>>(configJson);
        Console.WriteLine($"SDK_Init: Configuration fetched: {configJson}");

        IsInitialized = true;
    }

    public async Task<string> CreateSessionAsync(string identity, string secret)
    {
        _jwt = await _netClient.AuthenticateAsync(identity, secret, _apiEndpoint!);
        return _jwt;
    }

    public async Task KeyOperationAsync(string opType, string keyName, byte[]? keyData)
    {
        if (string.IsNullOrEmpty(_jwt))
        {
            throw new SDKException("No active session. Please call CreateSessionAsync() first.");
        }

        string response = await _netClient.KeyOpAsync(_jwt, opType, keyName, keyData, _apiEndpoint!);

        if ("READ".Equals(opType, StringComparison.OrdinalIgnoreCase))
        {
            // Assuming the response is the raw key material
            byte[] rawKey = Convert.FromHexString(response); // Use Hex in .NET 6+
            _keyCache.TryAdd(keyName, rawKey);
        }
        else if ("DELETE".Equals(opType, StringComparison.OrdinalIgnoreCase))
        {
            _keyCache.TryRemove(keyName, out _);
        }
    }

    public async Task<byte[]> GetKeyAsync(string keyName)
    {
        // Check cache first
        if (_keyCache.TryGetValue(keyName, out byte[]? key))
        {
            return key;
        }

        // If not in cache, fetch from backend
        Console.WriteLine($"Key '{keyName}' not in cache. Fetching from backend...");
        await KeyOperationAsync("READ", keyName, null);

        // Try getting from cache again
        if (_keyCache.TryGetValue(keyName, out byte[]? fetchedKey))
        {
            return fetchedKey;
        }

        throw new SDKException($"Key not found: {keyName}");
    }

    public void Cleanup()
    {
        // Securely clear sensitive data
        _jwt = null;
        foreach (var entry in _keyCache)
        {
            Array.Clear(entry.Value, 0, entry.Value.Length);
        }
        _keyCache.Clear();
        IsInitialized = false;
    }

    public void Dispose()
    {
        Cleanup();
        _netClient.Dispose();
    }
}
