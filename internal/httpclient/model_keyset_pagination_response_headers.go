/*
Ory Hydra API

Documentation for all of Ory Hydra's APIs.

API version:
Contact: hi@ory.sh
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the KeysetPaginationResponseHeaders type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &KeysetPaginationResponseHeaders{}

// KeysetPaginationResponseHeaders The `Link` HTTP header contains multiple links (`first`, `next`) formatted as: `<https://{project-slug}.projects.oryapis.com/admin/sessions?page_size=250&page_token=>; rel=\"first\"`  For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
type KeysetPaginationResponseHeaders struct {
	// The Link HTTP Header  The `Link` header contains a comma-delimited list of links to the following pages:  first: The first page of results. next: The next page of results.  Pages are omitted if they do not exist. For example, if there is no next page, the `next` link is omitted. Examples:  </admin/sessions?page_size=250&page_token={last_item_uuid}; rel=\"first\",/admin/sessions?page_size=250&page_token=>; rel=\"next\"
	Link *string `json:"link,omitempty"`
}

// NewKeysetPaginationResponseHeaders instantiates a new KeysetPaginationResponseHeaders object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKeysetPaginationResponseHeaders() *KeysetPaginationResponseHeaders {
	this := KeysetPaginationResponseHeaders{}
	return &this
}

// NewKeysetPaginationResponseHeadersWithDefaults instantiates a new KeysetPaginationResponseHeaders object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKeysetPaginationResponseHeadersWithDefaults() *KeysetPaginationResponseHeaders {
	this := KeysetPaginationResponseHeaders{}
	return &this
}

// GetLink returns the Link field value if set, zero value otherwise.
func (o *KeysetPaginationResponseHeaders) GetLink() string {
	if o == nil || IsNil(o.Link) {
		var ret string
		return ret
	}
	return *o.Link
}

// GetLinkOk returns a tuple with the Link field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KeysetPaginationResponseHeaders) GetLinkOk() (*string, bool) {
	if o == nil || IsNil(o.Link) {
		return nil, false
	}
	return o.Link, true
}

// HasLink returns a boolean if a field has been set.
func (o *KeysetPaginationResponseHeaders) HasLink() bool {
	if o != nil && !IsNil(o.Link) {
		return true
	}

	return false
}

// SetLink gets a reference to the given string and assigns it to the Link field.
func (o *KeysetPaginationResponseHeaders) SetLink(v string) {
	o.Link = &v
}

func (o KeysetPaginationResponseHeaders) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o KeysetPaginationResponseHeaders) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Link) {
		toSerialize["link"] = o.Link
	}
	return toSerialize, nil
}

type NullableKeysetPaginationResponseHeaders struct {
	value *KeysetPaginationResponseHeaders
	isSet bool
}

func (v NullableKeysetPaginationResponseHeaders) Get() *KeysetPaginationResponseHeaders {
	return v.value
}

func (v *NullableKeysetPaginationResponseHeaders) Set(val *KeysetPaginationResponseHeaders) {
	v.value = val
	v.isSet = true
}

func (v NullableKeysetPaginationResponseHeaders) IsSet() bool {
	return v.isSet
}

func (v *NullableKeysetPaginationResponseHeaders) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKeysetPaginationResponseHeaders(val *KeysetPaginationResponseHeaders) *NullableKeysetPaginationResponseHeaders {
	return &NullableKeysetPaginationResponseHeaders{value: val, isSet: true}
}

func (v NullableKeysetPaginationResponseHeaders) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKeysetPaginationResponseHeaders) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
