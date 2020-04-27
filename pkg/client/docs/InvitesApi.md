# \InvitesApi

All URIs are relative to *https://local.moov.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InvitesGet**](InvitesApi.md#InvitesGet) | **Get** /invites | List outstanding invites
[**InvitesInviteIDDelete**](InvitesApi.md#InvitesInviteIDDelete) | **Delete** /invites/{inviteID} | Delete an invite that was sent and invalidate the token.
[**InvitesPost**](InvitesApi.md#InvitesPost) | **Post** /invites | Send an email invite to a new user



## InvitesGet

> []Invite InvitesGet(ctx, optional)

List outstanding invites

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***InvitesGetOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a InvitesGetOpts struct


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


## InvitesInviteIDDelete

> InvitesInviteIDDelete(ctx, inviteID)

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


## InvitesPost

> Invite InvitesPost(ctx, invite)

Send an email invite to a new user

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**invite** | [**Invite**](Invite.md)|  | 

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

