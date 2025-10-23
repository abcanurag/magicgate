namespace MyCompany.CryptoSDK;

/// <summary>
/// Custom exception for all errors originating from the SDK.
/// </summary>
public class SDKException : Exception
{
    public SDKException(string message) : base(message) { }
    public SDKException(string message, Exception innerException) : base(message, innerException) { }
}
