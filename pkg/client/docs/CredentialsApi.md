# \CredentialsApi

All URIs are relative to *https://local.moov.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DisableCredentials**](CredentialsApi.md#DisableCredentials) | **Get** /identities/{identityID}/credentials/{credentialID} | Disables a credential so it can&#39;t be used anymore to login
[**ListCredentials**](CredentialsApi.md#ListCredentials) | **Get** /identities/{identityID}/credentials | List the credentials this user has used.



## DisableCredentials

> []Credential DisableCredentials(ctx, identityID, credentialID)

Disables a credential so it can't be used anymore to login

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**identityID** | [**string**](.md)| ID of the Identity for the credential | 
**credentialID** | [**string**](.md)| ID of the credential to disable | 

### Return type

[**[]Credential**](Credential.md)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCredentials

> []Credential ListCredentials(ctx, identityID)

List the credentials this user has used.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**identityID** | [**string**](.md)| ID of the Identity to lookup | 

### Return type

[**[]Credential**](Credential.md)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

