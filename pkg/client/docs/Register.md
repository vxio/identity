# Register

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Provider** | **string** | OIDC provider that was used to handle authentication of this user. | [optional] [readonly] 
**SubjectID** | **string** | ID of the remote OIDC server gives to this identity | [optional] [readonly] 
**InviteCode** | **string** |  | [optional] 
**FirstName** | **string** |  | [optional] 
**MiddleName** | **string** |  | [optional] 
**LastName** | **string** |  | [optional] 
**NickName** | Pointer to **string** |  | [optional] 
**Suffix** | Pointer to **string** |  | [optional] 
**BirthDate** | [**time.Time**](time.Time.md) |  | [optional] 
**Email** | **string** | Email Address | [optional] 
**Phones** | [**[]RegisterPhone**](RegisterPhone.md) |  | [optional] 
**Addresses** | [**[]RegisterAddress**](RegisterAddress.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


