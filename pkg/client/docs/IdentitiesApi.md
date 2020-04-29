# \IdentitiesApi

All URIs are relative to *https://local.moov.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DisableIdentity**](IdentitiesApi.md#DisableIdentity) | **Delete** /identities/{identityID} | Disable an identity. Its left around for historical reporting
[**GetIdentity**](IdentitiesApi.md#GetIdentity) | **Get** /identities/{identityID} | List identities and associates userId
[**ListIdentities**](IdentitiesApi.md#ListIdentities) | **Get** /identities | List identities and associates userId
[**UpdateIdentity**](IdentitiesApi.md#UpdateIdentity) | **Put** /identities/{identityID} | Update a specific Identity



## DisableIdentity

> DisableIdentity(ctx, identityID)

Disable an identity. Its left around for historical reporting

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**identityID** | [**string**](.md)| ID of the Identity to lookup | 

### Return type

 (empty response body)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIdentity

> Identity GetIdentity(ctx, identityID)

List identities and associates userId

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**identityID** | [**string**](.md)| ID of the Identity to lookup | 

### Return type

[**Identity**](Identity.md)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListIdentities

> []Identity ListIdentities(ctx, optional)

List identities and associates userId

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ListIdentitiesOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ListIdentitiesOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **orgID** | [**optional.Interface of string**](.md)| Filter in only for specific Organization | 

### Return type

[**[]Identity**](Identity.md)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateIdentity

> Identity UpdateIdentity(ctx, identityID, updateIdentity)

Update a specific Identity

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**identityID** | [**string**](.md)| ID of the Identity to lookup | 
**updateIdentity** | [**UpdateIdentity**](UpdateIdentity.md)|  | 

### Return type

[**Identity**](Identity.md)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

