# \InternalApi

All URIs are relative to *https://local.moov.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Authenticated**](InternalApi.md#Authenticated) | **Get** /authenticated | Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service redirect to this endpoint. 
[**Register**](InternalApi.md#Register) | **Get** /register | Returns the partially completed registration details that were pulled by AuthN service. 
[**RegisterWithCredentials**](InternalApi.md#RegisterWithCredentials) | **Post** /register | Called when the user is registering for the first time. It requires that they have authenticated with a  supported OIDC provider and recieved a valid invite code. 



## Authenticated

> LoggedIn Authenticated(ctx, )

Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service redirect to this endpoint. 

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**LoggedIn**](LoggedIn.md)

### Authorization

[LoginAuth](../README.md#LoginAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Register

> Register Register(ctx, )

Returns the partially completed registration details that were pulled by AuthN service. 

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**Register**](Register.md)

### Authorization

[LoginAuth](../README.md#LoginAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RegisterWithCredentials

> LoggedIn RegisterWithCredentials(ctx, register)

Called when the user is registering for the first time. It requires that they have authenticated with a  supported OIDC provider and recieved a valid invite code. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
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

