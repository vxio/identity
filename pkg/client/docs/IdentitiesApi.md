# \IdentitiesApi

All URIs are relative to *https://local.moov.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**IdentitiesGet**](IdentitiesApi.md#IdentitiesGet) | **Get** /identities | List identities and associates userId
[**IdentitiesIdentityIDCredentialsGet**](IdentitiesApi.md#IdentitiesIdentityIDCredentialsGet) | **Get** /identities/{identityID}/credentials | List the credentials this user has used.
[**IdentitiesIdentityIDGet**](IdentitiesApi.md#IdentitiesIdentityIDGet) | **Get** /identities/{identityID} | List identities and associates userId
[**IdentitiesIdentityIDPut**](IdentitiesApi.md#IdentitiesIdentityIDPut) | **Put** /identities/{identityID} | Update a specific Identity



## IdentitiesGet

> []Identity IdentitiesGet(ctx, optional)

List identities and associates userId

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***IdentitiesGetOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a IdentitiesGetOpts struct


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


## IdentitiesIdentityIDCredentialsGet

> []Credential IdentitiesIdentityIDCredentialsGet(ctx, identityID)

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


## IdentitiesIdentityIDGet

> Identity IdentitiesIdentityIDGet(ctx, identityID)

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


## IdentitiesIdentityIDPut

> Identity IdentitiesIdentityIDPut(ctx, identityID, identity)

Update a specific Identity

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**identityID** | [**string**](.md)| ID of the Identity to lookup | 
**identity** | [**Identity**](Identity.md)|  | 

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

