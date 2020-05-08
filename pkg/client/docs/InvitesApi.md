# \InvitesApi

All URIs are relative to *https://local.moov.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DisableInvite**](InvitesApi.md#DisableInvite) | **Delete** /invites/{inviteID} | Delete an invite that was sent and invalidate the token.
[**ListInvites**](InvitesApi.md#ListInvites) | **Get** /invites | List outstanding invites
[**SendInvite**](InvitesApi.md#SendInvite) | **Post** /invites | Send an email invite to a new user



## DisableInvite

> DisableInvite(ctx, inviteID)

Delete an invite that was sent and invalidate the token.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**inviteID** | [**string**](.md)| ID of the invite to delete | 

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


## ListInvites

> []Invite ListInvites(ctx, optional)

List outstanding invites

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ListInvitesOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ListInvitesOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **orgID** | [**optional.Interface of string**](.md)| Filter in only for specific Organization | 

### Return type

[**[]Invite**](Invite.md)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SendInvite

> Invite SendInvite(ctx, sendInvite)

Send an email invite to a new user

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**sendInvite** | [**SendInvite**](SendInvite.md)|  | 

### Return type

[**Invite**](Invite.md)

### Authorization

[GatewayAuth](../README.md#GatewayAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

