# KeysetPaginationResponseHeaders

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Link** | Pointer to **string** | The Link HTTP Header  The &#x60;Link&#x60; header contains a comma-delimited list of links to the following pages:  first: The first page of results. next: The next page of results.  Pages are omitted if they do not exist. For example, if there is no next page, the &#x60;next&#x60; link is omitted. Examples:  &lt;/admin/sessions?page_size&#x3D;250&amp;page_token&#x3D;{last_item_uuid}; rel&#x3D;\&quot;first\&quot;,/admin/sessions?page_size&#x3D;250&amp;page_token&#x3D;&gt;; rel&#x3D;\&quot;next\&quot; | [optional] 

## Methods

### NewKeysetPaginationResponseHeaders

`func NewKeysetPaginationResponseHeaders() *KeysetPaginationResponseHeaders`

NewKeysetPaginationResponseHeaders instantiates a new KeysetPaginationResponseHeaders object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeysetPaginationResponseHeadersWithDefaults

`func NewKeysetPaginationResponseHeadersWithDefaults() *KeysetPaginationResponseHeaders`

NewKeysetPaginationResponseHeadersWithDefaults instantiates a new KeysetPaginationResponseHeaders object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLink

`func (o *KeysetPaginationResponseHeaders) GetLink() string`

GetLink returns the Link field if non-nil, zero value otherwise.

### GetLinkOk

`func (o *KeysetPaginationResponseHeaders) GetLinkOk() (*string, bool)`

GetLinkOk returns a tuple with the Link field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLink

`func (o *KeysetPaginationResponseHeaders) SetLink(v string)`

SetLink sets Link field to given value.

### HasLink

`func (o *KeysetPaginationResponseHeaders) HasLink() bool`

HasLink returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


