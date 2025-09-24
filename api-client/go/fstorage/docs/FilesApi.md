# \FilesApi

All URIs are relative to *http://localhost:3000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1FilesBatchDelete**](FilesApi.md#ApiV1FilesBatchDelete) | **Delete** /api/v1/files/batch | Delete multiple files
[**ApiV1FilesUploadPost**](FilesApi.md#ApiV1FilesUploadPost) | **Post** /api/v1/files/upload | Upload multiple files


# **ApiV1FilesBatchDelete**
> ApiV1FilesBatchDelete(ctx, body)
Delete multiple files

Delete multiple files by their IDs from UploadThing storage

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**DeleteRequest**](DeleteRequest.md)| File IDs to delete | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1FilesUploadPost**
> InlineResponse200 ApiV1FilesUploadPost(ctx, files, domain)
Upload multiple files

Upload one or more files to UploadThing storage with a specified domain

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **files** | [**[]*os.File**](*os.File.md)| Files to upload | 
  **domain** | **string**| Domain for all uploaded files | 

### Return type

[**InlineResponse200**](inline_response_200.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

