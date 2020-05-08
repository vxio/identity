# \InternalApi

All URIs are relative to *https://local.moov.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**LoginWithCredentials**](InternalApi.md#LoginWithCredentials) | **Post** /login | Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service will call  this endpoint to record and finish the login to get their token to use the API.  If the client service recieves a 404 they must send them to registration if its allowed per the client or check for an invite for authenticated users email before sending to registration.       
[**RegisterWithCredentials**](InternalApi.md#RegisterWithCredentials) | **Post** /register | Register user based on OIDC credentials.  This is called by the OIDC client services we create to register the user with what  available information they have and obtain from the user. 



## LoginWithCredentials

> LoggedIn LoginWithCredentials(ctx, login)

Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service will call  this endpoint to record and finish the login to get their token to use the API.  If the client service recieves a 404 they must send them to registration if its allowed per the client or check for an invite for authenticated users email before sending to registration.       

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**login** | [**Login**](Login.md)| Arguments needed to match up the OIDC user data with a user in the system | 

### Return type

[**LoggedIn**](LoggedIn.md)

### Authorization

[ServiceAuth](../README.md#ServiceAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RegisterWithCredentials

> LoggedIn RegisterWithCredentials(ctx, register)

Register user based on OIDC credentials.  This is called by the OIDC client services we create to register the user with what  available information they have and obtain from the user. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**register** | [**Register**](Register.md)| Arguments needed register a user with OIDC credentials. | 

### Return type

[**LoggedIn**](LoggedIn.md)

### Authorization

[ServiceAuth](../README.md#ServiceAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

