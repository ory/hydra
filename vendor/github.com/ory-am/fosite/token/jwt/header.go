package jwt

import "github.com/dgrijalva/jwt-go"

// Headers is the jwt headers
type Headers struct {
	Extra map[string]interface{}
}

// ToMap will transform the headers to a map structure
func (h *Headers) ToMap() map[string]interface{} {
	var filter = map[string]bool{"alg": true, "typ": true}
	var extra = map[string]interface{}{}

	// filter known values from extra.
	for k, v := range h.Extra {
		if _, ok := filter[k]; !ok {
			extra[k] = v
		}
	}

	return extra
}

// Add will add a key-value pair to the extra field
func (h *Headers) Add(key string, value interface{}) {
	if h.Extra == nil {
		h.Extra = make(map[string]interface{})
	}
	h.Extra[key] = value
}

// Get will get a value from the extra field based on a given key
func (h *Headers) Get(key string) interface{} {
	return h.Extra[key]
}

// ToMapClaims will return a jwt-go MapClaims representaion
func (h Headers) ToMapClaims() jwt.MapClaims {
	return h.ToMap()
}
