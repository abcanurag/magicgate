using System.Text;
using MyCompany.CryptoSDK;

Console.WriteLine("--- .NET SDK Usage Example ---\n");

try
{
    // 1. Initialize the SDK
    Console.WriteLine("1. Initializing SDK...");
    await CryptoSDK.InitAsync("my-app-registration-token-12345");
    Console.WriteLine();

    // 2. Create a session
    Console.WriteLine("2. Creating a session...");
    string jwt = await CryptoSDK.CreateSessionAsync("app_user_01", "super_secret_password");
    Console.WriteLine($"   Session created. Received JWT: {jwt.Substring(0, 30)}...\n");

    // 3. Manage a key (CREATE)
    Console.WriteLine("3. Creating a key named 'MySecretKey'...");
    // In a real scenario, the key material might be generated locally and sent to the backend.
    byte[] keyMaterial = Encoding.UTF8.GetBytes("this-is-my-super-secret-key-data");
    await CryptoSDK.KeyOperationAsync("CREATE", "MySecretKey", keyMaterial);
    Console.WriteLine("   Key 'MySecretKey' created on backend.\n");

    // 4. Perform crypto operations (Encrypt and Decrypt)
    Console.WriteLine("4. Encrypting and decrypting data with 'MySecretKey'...");
    string plaintext = "This is a very sensitive message.";
    string algo = "AES-256-GCM";

    // Encrypt
    byte[] ciphertext = await CryptoSDK.DoCryptoAsync("ENCRYPT", "MySecretKey", algo, Encoding.UTF8.GetBytes(plaintext));
    Console.WriteLine("   Encryption successful.");
    Console.WriteLine($"   Plaintext: '{plaintext}'");
    Console.WriteLine($"   Ciphertext (hex): {Convert.ToHexString(ciphertext)}\n");

    // Decrypt
    byte[] decryptedTextBytes = await CryptoSDK.DoCryptoAsync("DECRYPT", "MySecretKey", algo, ciphertext);
    string decryptedText = Encoding.UTF8.GetString(decryptedTextBytes);
    Console.WriteLine("   Decryption successful.");
    Console.WriteLine($"   Decrypted Text: '{decryptedText}'\n");

    // Verify
    if (plaintext != decryptedText) throw new Exception("FATAL: Decrypted text does not match original plaintext!");
    Console.WriteLine("   SUCCESS: Original plaintext matches decrypted text.\n");

    // 5. Demonstrate key caching
    Console.WriteLine("5. Encrypting again (should use cached key)...");
    byte[] ciphertext2 = await CryptoSDK.DoCryptoAsync("ENCRYPT", "MySecretKey", algo, Encoding.UTF8.GetBytes(plaintext));
    if (ciphertext.SequenceEqual(ciphertext2)) throw new Exception("FATAL: Ciphertext should be different on each encryption!");
    Console.WriteLine("   Second encryption successful (and produced different ciphertext due to random nonce).\n");

}
catch (SDKException e)
{
    Console.Error.WriteLine($"SDK operation failed: {e.Message}");
    if (e.InnerException != null) Console.Error.WriteLine($"  -> Inner Exception: {e.InnerException.Message}");
}
finally
{
    // 6. Clean up
    Console.WriteLine("6. Cleaning up SDK resources...");
    CryptoSDK.Cleanup();
}

Console.WriteLine("\n--- SDK Example Finished ---");
