# \InternalApi

All URIs are relative to *https://local.moov.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Authenticated**](InternalApi.md#Authenticated) | **Post** /authenticated | Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service redirect to this endpoint.     
[**RegisterWithCredentials**](InternalApi.md#RegisterWithCredentials) | **Post** /register | If the OIDC client specified it got an invite code that token will be exchanged here to login 



## Authenticated

> LoggedIn Authenticated(ctx, moovLogin, login)

Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service redirect to this endpoint.     

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**moovLogin** | **string**| Encrypted and signed token that they authenticated via one of the approved services | 
**login** | [**Login**](Login.md)| Arguments needed to match up the OIDC user data with a user in the system | 

### Return type

[**LoggedIn**](LoggedIn.md)

### Authorization

[LoginAuth](../README.md#LoginAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RegisterWithCredentials

> LoggedIn RegisterWithCredentials(ctx, moovLogin, register)

If the OIDC client specified it got an invite code that token will be exchanged here to login 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**moovLogin** | **string**| Encrypted and signed token that they authenticated via one of the approved services | 
**register** | [**Register**](Register.md)| Arguments needed register a user with OIDC credentials. | 

### Return type

[**LoggedIn**](LoggedIn.md)

### Authorization

[LoginAuth](../README.md#LoginAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

