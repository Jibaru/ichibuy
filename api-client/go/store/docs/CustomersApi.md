# \CustomersApi

All URIs are relative to *https://ichibuy-store.vercel.app*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1CustomersIdDelete**](CustomersApi.md#ApiV1CustomersIdDelete) | **Delete** /api/v1/customers/{id} | Delete customer by ID
[**ApiV1CustomersIdGet**](CustomersApi.md#ApiV1CustomersIdGet) | **Get** /api/v1/customers/{id} | Get customer by ID
[**ApiV1CustomersIdPut**](CustomersApi.md#ApiV1CustomersIdPut) | **Put** /api/v1/customers/{id} | Update customer by ID
[**ApiV1CustomersPost**](CustomersApi.md#ApiV1CustomersPost) | **Post** /api/v1/customers | Create a new customer


# **ApiV1CustomersIdDelete**
> ApiV1CustomersIdDelete(ctx, id)
Delete customer by ID

Delete a customer

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| Customer ID | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1CustomersIdGet**
> ServicesGetCustomerResp ApiV1CustomersIdGet(ctx, id)
Get customer by ID

Retrieve a customer by its ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| Customer ID | 

### Return type

[**ServicesGetCustomerResp**](services.GetCustomerResp.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1CustomersIdPut**
> ApiV1CustomersIdPut(ctx, id, customer)
Update customer by ID

Update a customer's information

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| Customer ID | 
  **customer** | [**HandlersUpdateCustomerBody**](HandlersUpdateCustomerBody.md)| Customer data | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1CustomersPost**
> ServicesCreateCustomerResp ApiV1CustomersPost(ctx, customer)
Create a new customer

Create a new customer with first name, last name, email and phone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **customer** | [**HandlersCreateCustomerBody**](HandlersCreateCustomerBody.md)| Customer data | 

### Return type

[**ServicesCreateCustomerResp**](services.CreateCustomerResp.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

