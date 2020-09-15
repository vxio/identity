# \SessionApi

All URIs are relative to *https://local.moov.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ChangeSessionDetails**](SessionApi.md#ChangeSessionDetails) | **Put** /session | Changes the details of the session allowing to change tenants or identities
[**GetSessionDetails**](SessionApi.md#GetSessionDetails) | **Get** /session | Return information about the current session



## ChangeSessionDetails

> SessionDetails ChangeSessionDetails(ctx, changeSessionDetails)

Changes the details of the session allowing to change tenants or identities

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**changeSessionDetails** | [**ChangeSessionDetails**](ChangeSessionDetails.md)|  | 

### Return type

[**SessionDetails**](SessionDetails.md)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSessionDetails

> SessionDetails GetSessionDetails(ctx, )

Return information about the current session

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**SessionDetails**](SessionDetails.md)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

