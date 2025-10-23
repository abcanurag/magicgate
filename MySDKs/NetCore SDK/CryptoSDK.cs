using System.Security.Cryptography;
using MyCompany.CryptoSDK.Internal;

namespace MyCompany.CryptoSDK;

/// <summary>
/// Main public interface for the Crypto SDK.
/// This class provides static methods to initialize the SDK, manage sessions,
/// handle keys, and perform cryptographic operations.
/// </summary>
public static class CryptoSDK
{
    private static SDKContext? _context;
    private static readonly object _lock = new();

    /// <summary>
    /// Initializes the SDK. Must be called once before any other operation.
    /// </summary>
    /// <param name="regToken">The registration token for the application.</param>
    /// <exception cref="SDKException">Thrown if initialization fails or if already initialized.</exception>
    public static async Task InitAsync(string regToken)
    {
        lock (_lock)
        {
            if (_context?.IsInitialized == true)
            {
                throw new SDKException("SDK is already initialized.");
            }
            _context = new SDKContext();
        }
        await _context.InitializeAsync(regToken);
        Console.WriteLine("SDK Initialized successfully.");
    }

    /// <summary>
    /// Cleans up all resources used by the SDK.
    /// </summary>
    public static void Cleanup()
    {
        lock (_lock)
        {
            if (_context != null)
            {
                _context.Cleanup();
                _context = null;
                Console.WriteLine("SDK Cleanup complete.");
            }
        }
    }

    /// <summary>
    /// Creates an authenticated session with the backend service.
    /// </summary>
    /// <param name="identity">The user or application identity.</param>
    /// <param name="secret">The corresponding secret.</param>
    /// <returns>The received JWT for the session.</returns>
    /// <exception cref="SDKException">Thrown if the SDK is not initialized or authentication fails.</exception>
    public static async Task<string> CreateSessionAsync(string identity, string secret)
    {
        var ctx = CheckInitialized();
        return await ctx.CreateSessionAsync(identity, secret);
    }

    /// <summary>
    /// Performs a key management operation (CRUD).
    /// </summary>
    /// <param name="opType">The operation: "CREATE", "READ", "UPDATE", "DELETE".</param>
    /// <param name="keyName">The unique name of the key.</param>
    /// <param name="keyData">For "CREATE" or "UPDATE", the raw key material. Can be null.</param>
    /// <exception cref="SDKException">Thrown if the operation fails.</exception>
    public static async Task KeyOperationAsync(string opType, string keyName, byte[]? keyData)
    {
        var ctx = CheckInitialized();
        await ctx.KeyOperationAsync(opType, keyName, keyData);
    }

    /// <summary>
    /// Performs a cryptographic operation (encryption or decryption).
    /// </summary>
    /// <param name="opType">The crypto operation to perform ("ENCRYPT" or "DECRYPT").</param>
    /// <param name="keyName">The name of the key to use.</param>
    /// <param name="algoName">The algorithm to use (e.g., "AES-256-GCM").</param>
    /// <param name="input">The data to be processed (plaintext or ciphertext).</param>
    /// <returns>The result of the operation (ciphertext or plaintext).</returns>
    /// <exception cref="SDKException">Thrown if the operation fails.</exception>
    public static async Task<byte[]> DoCryptoAsync(string opType, string keyName, string algoName, byte[] input)
    {
        var ctx = CheckInitialized();

        // 1. Get key from context (fetches from backend if not in cache)
        byte[] key = await ctx.GetKeyAsync(keyName);

        // 2. Perform crypto operation
        if ("ENCRYPT".Equals(opType, StringComparison.OrdinalIgnoreCase))
        {
            return CryptoUtils.Encrypt(key, algoName, input);
        }
        if ("DECRYPT".Equals(opType, StringComparison.OrdinalIgnoreCase))
        {
            return CryptoUtils.Decrypt(key, algoName, input);
        }
        
        throw new SDKException($"Unsupported crypto operation: {opType}");
    }

    /// <summary>
    /// Helper to ensure the SDK is initialized before use.
    /// </summary>
    private static SDKContext CheckInitialized()
    {
        if (_context?.IsInitialized != true)
        {
            throw new SDKException("SDK is not initialized. Please call SDK.InitAsync() first.");
        }
        return _context;
    }
}
