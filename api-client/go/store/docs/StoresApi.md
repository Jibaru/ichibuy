# \StoresApi

All URIs are relative to *https://ichibuy-store.vercel.app*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1StoresGet**](StoresApi.md#ApiV1StoresGet) | **Get** /api/v1/stores | List stores
[**ApiV1StoresIdDelete**](StoresApi.md#ApiV1StoresIdDelete) | **Delete** /api/v1/stores/{id} | Delete store by ID
[**ApiV1StoresIdGet**](StoresApi.md#ApiV1StoresIdGet) | **Get** /api/v1/stores/{id} | Get store by ID
[**ApiV1StoresIdPut**](StoresApi.md#ApiV1StoresIdPut) | **Put** /api/v1/stores/{id} | Update store by ID
[**ApiV1StoresPost**](StoresApi.md#ApiV1StoresPost) | **Post** /api/v1/stores | Create a new store


# **ApiV1StoresGet**
> ServicesListStoresResp ApiV1StoresGet(ctx, optional)
List stores

Get paginated list of stores with filters and sorting

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***StoresApiApiV1StoresGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a StoresApiApiV1StoresGetOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **optional.String**| Filter by name | 
 **description** | **optional.String**| Filter by description | 
 **sortBy** | **optional.String**| Sort by field | [default to &quot;name&quot;]
 **sortOrder** | **optional.String**| Sort order | [default to &quot;ASC&quot;]
 **offset** | **optional.Int32**| Offset | [default to 0]
 **limit** | **optional.Int32**| Limit | [default to 10]

### Return type

[**ServicesListStoresResp**](services.ListStoresResp.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1StoresIdDelete**
> ApiV1StoresIdDelete(ctx, id)
Delete store by ID

Delete a store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| Store ID | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1StoresIdGet**
> ServicesGetStoreResp ApiV1StoresIdGet(ctx, id)
Get store by ID

Retrieve a store by its ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| Store ID | 

### Return type

[**ServicesGetStoreResp**](services.GetStoreResp.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1StoresIdPut**
> ApiV1StoresIdPut(ctx, id, store)
Update store by ID

Update a store's information

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| Store ID | 
  **store** | [**HandlersUpdateStoreBody**](HandlersUpdateStoreBody.md)| Store data | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1StoresPost**
> ServicesCreateStoreResp ApiV1StoresPost(ctx, store)
Create a new store

Create a new store with name, description and location

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **store** | [**HandlersCreateStoreBody**](HandlersCreateStoreBody.md)| Store data | 

### Return type

[**ServicesCreateStoreResp**](services.CreateStoreResp.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

