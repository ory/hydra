# TokenPaginationResponseHeaders

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Link** | Pointer to **string** | The Link HTTP Header  The &#x60;Link&#x60; header contains a comma-delimited list of links to the following pages:  first: The first page of results. next: The next page of results. prev: The previous page of results. last: The last page of results.  Pages are omitted if they do not exist. For example, if there is no next page, the &#x60;next&#x60; link is omitted. Examples:  &lt;/clients?page_size&#x3D;5&amp;page_token&#x3D;0&gt;; rel&#x3D;\&quot;first\&quot;,&lt;/clients?page_size&#x3D;5&amp;page_token&#x3D;15&gt;; rel&#x3D;\&quot;next\&quot;,&lt;/clients?page_size&#x3D;5&amp;page_token&#x3D;5&gt;; rel&#x3D;\&quot;prev\&quot;,&lt;/clients?page_size&#x3D;5&amp;page_token&#x3D;20&gt;; rel&#x3D;\&quot;last\&quot; | [optional] 
**XTotalCount** | Pointer to **int64** | The X-Total-Count HTTP Header  The &#x60;X-Total-Count&#x60; header contains the total number of items in the collection. | [optional] 

## Methods

### NewTokenPaginationResponseHeaders

`func NewTokenPaginationResponseHeaders() *TokenPaginationResponseHeaders`

NewTokenPaginationResponseHeaders instantiates a new TokenPaginationResponseHeaders object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTokenPaginationResponseHeadersWithDefaults

`func NewTokenPaginationResponseHeadersWithDefaults() *TokenPaginationResponseHeaders`

NewTokenPaginationResponseHeadersWithDefaults instantiates a new TokenPaginationResponseHeaders object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLink

`func (o *TokenPaginationResponseHeaders) GetLink() string`

GetLink returns the Link field if non-nil, zero value otherwise.

### GetLinkOk

`func (o *TokenPaginationResponseHeaders) GetLinkOk() (*string, bool)`

GetLinkOk returns a tuple with the Link field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLink

`func (o *TokenPaginationResponseHeaders) SetLink(v string)`

SetLink sets Link field to given value.

### HasLink

`func (o *TokenPaginationResponseHeaders) HasLink() bool`

HasLink returns a boolean if a field has been set.

### GetXTotalCount

`func (o *TokenPaginationResponseHeaders) GetXTotalCount() int64`

GetXTotalCount returns the XTotalCount field if non-nil, zero value otherwise.

### GetXTotalCountOk

`func (o *TokenPaginationResponseHeaders) GetXTotalCountOk() (*int64, bool)`

GetXTotalCountOk returns a tuple with the XTotalCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetXTotalCount

`func (o *TokenPaginationResponseHeaders) SetXTotalCount(v int64)`

SetXTotalCount sets XTotalCount field to given value.

### HasXTotalCount

`func (o *TokenPaginationResponseHeaders) HasXTotalCount() bool`

HasXTotalCount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


