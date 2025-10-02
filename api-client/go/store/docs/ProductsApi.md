# \ProductsApi

All URIs are relative to *https://ichibuy-store.vercel.app*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1ProductsGet**](ProductsApi.md#ApiV1ProductsGet) | **Get** /api/v1/products | List products
[**ApiV1ProductsIdDelete**](ProductsApi.md#ApiV1ProductsIdDelete) | **Delete** /api/v1/products/{id} | Delete product by ID
[**ApiV1ProductsIdGet**](ProductsApi.md#ApiV1ProductsIdGet) | **Get** /api/v1/products/{id} | Get a product by ID
[**ApiV1ProductsIdPut**](ProductsApi.md#ApiV1ProductsIdPut) | **Put** /api/v1/products/{id} | Update a product
[**ApiV1ProductsPost**](ProductsApi.md#ApiV1ProductsPost) | **Post** /api/v1/products | Create a new product


# **ApiV1ProductsGet**
> ServicesListProductsResp ApiV1ProductsGet(ctx, optional)
List products

Get paginated list of products with filters and sorting

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ProductsApiApiV1ProductsGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ProductsApiApiV1ProductsGetOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **storeId** | **optional.String**| Filter by store ID | 
 **name** | **optional.String**| Filter by name | 
 **description** | **optional.String**| Filter by description | 
 **active** | **optional.Bool**| Filter by active status | 
 **sortBy** | **optional.String**| Sort by field | [default to &quot;name&quot;]
 **sortOrder** | **optional.String**| Sort order | [default to &quot;ASC&quot;]
 **offset** | **optional.Int32**| Offset | [default to 0]
 **limit** | **optional.Int32**| Limit | [default to 10]

### Return type

[**ServicesListProductsResp**](services.ListProductsResp.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1ProductsIdDelete**
> ApiV1ProductsIdDelete(ctx, id)
Delete product by ID

Delete a product

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| Product ID | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1ProductsIdGet**
> ServicesGetProductResp ApiV1ProductsIdGet(ctx, id)
Get a product by ID

Get a single product by its ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| Product ID | 

### Return type

[**ServicesGetProductResp**](services.GetProductResp.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1ProductsIdPut**
> ApiV1ProductsIdPut(ctx, id, name, active, prices, optional)
Update a product

Update an existing product. Note: StoreID, ID and CreatedAt cannot be updated. Images and Prices are replaced entirely.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| Product ID | 
  **name** | **string**| Product name | 
  **active** | **bool**| Product active status | 
  **prices** | **string**| JSON array of prices | 
 **optional** | ***ProductsApiApiV1ProductsIdPutOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ProductsApiApiV1ProductsIdPutOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **description** | **optional.String**| Product description | 
 **deleteImagesIDs** | **optional.String**| JSON array of image IDs to delete | 
 **deletePricesIDs** | **optional.String**| JSON array of price IDs to delete | 
 **images** | **optional.Interface of *os.File**| Product images (multiple files allowed) | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1ProductsPost**
> ServicesCreateProductResp ApiV1ProductsPost(ctx, name, active, storeId, prices, optional)
Create a new product

Create a new product with name, description, active status, store ID, images (as files) and prices

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Product name | 
  **active** | **bool**| Product active status | 
  **storeId** | **string**| Store ID | 
  **prices** | **string**| JSON array of prices | 
 **optional** | ***ProductsApiApiV1ProductsPostOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ProductsApiApiV1ProductsPostOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **description** | **optional.String**| Product description | 
 **images** | **optional.Interface of *os.File**| Product images (multiple files allowed) | 

### Return type

[**ServicesCreateProductResp**](services.CreateProductResp.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

