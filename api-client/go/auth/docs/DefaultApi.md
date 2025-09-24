# \DefaultApi

All URIs are relative to *https://ichibuy-auth.vercel.app*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1AuthProviderCallbackGet**](DefaultApi.md#ApiV1AuthProviderCallbackGet) | **Get** /api/v1/auth/{provider}/callback | OAuthCallback
[**ApiV1AuthProviderGet**](DefaultApi.md#ApiV1AuthProviderGet) | **Get** /api/v1/auth/{provider} | StartOAuth
[**ApiV1AuthWellKnownJwksJsonGet**](DefaultApi.md#ApiV1AuthWellKnownJwksJsonGet) | **Get** /api/v1/auth/.well-known/jwks.json | GetJWKS


# **ApiV1AuthProviderCallbackGet**
> ApiV1AuthProviderCallbackGet(ctx, provider)
OAuthCallback

OAuthCallback

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **provider** | **string**| Provider: google | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1AuthProviderGet**
> ApiV1AuthProviderGet(ctx, provider)
StartOAuth

StartOAuth

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **provider** | **string**| Provider: google | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1AuthWellKnownJwksJsonGet**
> HandlersJwks ApiV1AuthWellKnownJwksJsonGet(ctx, )
GetJWKS

Returns JSON Web Key Set for JWT verification

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**HandlersJwks**](handlers.JWKS.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

