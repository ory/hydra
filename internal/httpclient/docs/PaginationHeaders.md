# PaginationHeaders

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Link** | Pointer to **string** | The link header contains pagination links.  For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).  in: header | [optional] 
**XTotalCount** | Pointer to **string** | The total number of clients.  in: header | [optional] 

## Methods

### NewPaginationHeaders

`func NewPaginationHeaders() *PaginationHeaders`

NewPaginationHeaders instantiates a new PaginationHeaders object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPaginationHeadersWithDefaults

`func NewPaginationHeadersWithDefaults() *PaginationHeaders`

NewPaginationHeadersWithDefaults instantiates a new PaginationHeaders object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLink

`func (o *PaginationHeaders) GetLink() string`

GetLink returns the Link field if non-nil, zero value otherwise.

### GetLinkOk

`func (o *PaginationHeaders) GetLinkOk() (*string, bool)`

GetLinkOk returns a tuple with the Link field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLink

`func (o *PaginationHeaders) SetLink(v string)`

SetLink sets Link field to given value.

### HasLink

`func (o *PaginationHeaders) HasLink() bool`

HasLink returns a boolean if a field has been set.

### GetXTotalCount

`func (o *PaginationHeaders) GetXTotalCount() string`

GetXTotalCount returns the XTotalCount field if non-nil, zero value otherwise.

### GetXTotalCountOk

`func (o *PaginationHeaders) GetXTotalCountOk() (*string, bool)`

GetXTotalCountOk returns a tuple with the XTotalCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetXTotalCount

`func (o *PaginationHeaders) SetXTotalCount(v string)`

SetXTotalCount sets XTotalCount field to given value.

### HasXTotalCount

`func (o *PaginationHeaders) HasXTotalCount() bool`

HasXTotalCount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


