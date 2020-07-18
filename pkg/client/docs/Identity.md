# Identity

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IdentityID** | **string** | UUID v4 | [optional] 
**TenantID** | **string** | UUID v4 | [optional] 
**FirstName** | **string** |  | 
**MiddleName** | **string** |  | [optional] 
**LastName** | **string** |  | 
**NickName** | Pointer to **string** |  | [optional] 
**Suffix** | Pointer to **string** |  | [optional] 
**BirthDate** | [**time.Time**](time.Time.md) |  | [optional] 
**Status** | **string** |  | [optional] 
**Email** | **string** | Email Address | 
**EmailVerified** | **bool** | The user has verified they have access to this email | [optional] [readonly] 
**Phones** | [**[]Phone**](Phone.md) |  | [optional] 
**Addresses** | [**[]Address**](Address.md) |  | [optional] 
**RegisteredOn** | [**time.Time**](time.Time.md) |  | [optional] 
**LastLogin** | [**LastLogin**](LastLogin.md) |  | [optional] 
**DisabledOn** | Pointer to [**time.Time**](time.Time.md) |  | [optional] 
**DisabledBy** | Pointer to **string** | UUID v4 | [optional] 
**LastUpdatedOn** | [**time.Time**](time.Time.md) |  | [optional] 
**InviteID** | Pointer to **string** | UUID v4 | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


