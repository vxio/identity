# Credential

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CredentialID** | **string** | UUID v4 | [optional] 
**Provider** | **string** | OIDC provider that was used to handle authentication of this user. | [optional] [readonly] 
**SubjectID** | **string** | ID of the remote OIDC server gives to this identity | [optional] [readonly] 
**IdentityID** | **string** | UUID v4 | [optional] 
**TenantID** | **string** | UUID v4 | [optional] 
**CreatedOn** | [**time.Time**](time.Time.md) |  | [optional] 
**LastUsedOn** | [**time.Time**](time.Time.md) |  | [optional] 
**DisabledOn** | Pointer to [**time.Time**](time.Time.md) |  | [optional] 
**DisabledBy** | Pointer to **string** | UUID v4 | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


