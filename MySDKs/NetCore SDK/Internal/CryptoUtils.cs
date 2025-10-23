using System.Security.Cryptography;

namespace MyCompany.CryptoSDK.Internal;

/// <summary>
/// Utility class for cryptographic operations using .NET libraries.
/// </summary>
internal static class CryptoUtils
{
    private const int GCM_NONCE_SIZE = 12; // 96 bits
    private const int GCM_TAG_SIZE = 16;   // 128 bits

    public static byte[] Encrypt(byte[] key, string algoName, byte[] plaintext)
    {
        if (!"AES-256-GCM".Equals(algoName, StringComparison.OrdinalIgnoreCase))
            throw new SDKException($"Unsupported algorithm: {algoName}");

        // Nonce must be unique for each encryption with the same key
        byte[] nonce = new byte[GCM_NONCE_SIZE];
        RandomNumberGenerator.Fill(nonce);

        byte[] ciphertext = new byte[plaintext.Length];
        byte[] tag = new byte[GCM_TAG_SIZE];

        using var aesGcm = new AesGcm(key);
        aesGcm.Encrypt(nonce, plaintext, ciphertext, tag);

        // Prepend nonce and append tag to the ciphertext
        byte[] encryptedData = new byte[nonce.Length + ciphertext.Length + tag.Length];
        Buffer.BlockCopy(nonce, 0, encryptedData, 0, nonce.Length);
        Buffer.BlockCopy(ciphertext, 0, encryptedData, nonce.Length, ciphertext.Length);
        Buffer.BlockCopy(tag, 0, encryptedData, nonce.Length + ciphertext.Length, tag.Length);

        return encryptedData;
    }

    public static byte[] Decrypt(byte[] key, string algoName, byte[] encryptedData)
    {
        if (!"AES-256-GCM".Equals(algoName, StringComparison.OrdinalIgnoreCase))
            throw new SDKException($"Unsupported algorithm: {algoName}");

        try
        {
            ReadOnlySpan<byte> nonce = encryptedData.AsSpan(0, GCM_NONCE_SIZE);
            ReadOnlySpan<byte> tag = encryptedData.AsSpan(encryptedData.Length - GCM_TAG_SIZE);
            ReadOnlySpan<byte> ciphertext = encryptedData.AsSpan(GCM_NONCE_SIZE, encryptedData.Length - GCM_NONCE_SIZE - GCM_TAG_SIZE);

            byte[] decryptedData = new byte[ciphertext.Length];

            using var aesGcm = new AesGcm(key);
            aesGcm.Decrypt(nonce, ciphertext, tag, decryptedData);

            return decryptedData;
        }
        catch (CryptographicException e)
        {
            throw new SDKException("Decryption failed. The data may be corrupt, the key incorrect, or the tag invalid.", e);
        }
    }
}
