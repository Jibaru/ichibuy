# \GraphqlApi

All URIs are relative to *https://ichibuy-store.vercel.app*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1GraphqlPost**](GraphqlApi.md#ApiV1GraphqlPost) | **Post** /api/v1/graphql | GraphQL endpoint for stores


# **ApiV1GraphqlPost**
> interface{} ApiV1GraphqlPost(ctx, query)
GraphQL endpoint for stores

GraphQL endpoint to query stores with filters, sorting and pagination

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **query** | [**Query**](Query.md)| GraphQL query | 

### Return type

**interface{}**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

