// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package decoderx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/ory/x/assertx"

	"github.com/tidwall/gjson"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/jsonschema/v3"
)

func newRequest(t *testing.T, method, url string, body io.Reader, ct string) *http.Request {
	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", ct)
	return req
}

func TestHTTPFormDecoder(t *testing.T) {
	for k, tc := range []struct {
		d             string
		request       *http.Request
		contentType   string
		options       []HTTPDecoderOption
		expected      string
		expectedError string
	}{
		{
			d:             "should fail because the method is GET",
			request:       &http.Request{Header: map[string][]string{}, Method: "GET"},
			expectedError: "HTTP Request Method",
		},
		{
			d:             "should fail because the body is empty",
			request:       &http.Request{Header: map[string][]string{}, Method: "POST"},
			expectedError: "Content-Length",
		},
		{
			d:             "should fail because content type is missing",
			request:       newRequest(t, "POST", "/", nil, ""),
			expectedError: "Content-Length",
		},
		{
			d:             "should fail because content type is missing",
			request:       newRequest(t, "POST", "/", bytes.NewBufferString("foo"), ""),
			expectedError: "Content-Type",
		},
		{
			d:        "should pass with json without validation",
			request:  newRequest(t, "POST", "/", bytes.NewBufferString(`{"foo":"bar"}`), httpContentTypeJSON),
			expected: `{"foo":"bar"}`,
		},
		{
			d:             "should fail json if content type is not accepted",
			request:       newRequest(t, "POST", "/", bytes.NewBufferString(`{"foo":"bar"}`), httpContentTypeJSON),
			options:       []HTTPDecoderOption{HTTPFormDecoder()},
			expectedError: "Content-Type: application/json",
		},
		{
			d:       "should fail json if validation fails",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(`{"foo":"bar", "bar":"baz"}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{HTTPJSONDecoder(), MustHTTPRawJSONSchemaCompiler([]byte(`{
	"$id": "https://example.com/config.schema.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
		"foo": {
			"type": "number"
		},
		"bar": {
			"type": "string"
		}
	}
}`),
			)},
			expectedError: "expected number, but got string",
			expected:      `{ "bar": "baz", "foo": "bar" }`,
		},
		{
			d:       "should pass json with validation",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(`{"foo":"bar"}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{HTTPJSONDecoder(), MustHTTPRawJSONSchemaCompiler([]byte(`{
	"$id": "https://example.com/config.schema.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
		"foo": {
			"type": "string"
		}
	}
}`),
			),
			},
			expected: `{"foo":"bar"}`,
		},
		{
			d:             "should fail form request when form is used but only json is allowed",
			request:       newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{"foo": {"bar"}}.Encode()), httpContentTypeURLEncodedForm),
			options:       []HTTPDecoderOption{HTTPJSONDecoder()},
			expectedError: "Content-Type: application/x-www-form-urlencoded",
		},
		{
			d:             "should fail form request when schema is missing",
			request:       newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{"foo": {"bar"}}.Encode()), httpContentTypeURLEncodedForm),
			options:       []HTTPDecoderOption{},
			expectedError: "no validation schema was provided",
		},
		{
			d:             "should fail form request when schema does not validate request",
			request:       newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{"bar": {"bar"}}.Encode()), httpContentTypeURLEncodedForm),
			options:       []HTTPDecoderOption{HTTPJSONSchemaCompiler("stub/schema.json", nil)},
			expectedError: `missing properties: "foo"`,
		},
		{
			d:       "should fail for invalid JSON data with unrestricted object",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(`{"dynamic_object":{"stuff":{"blub":[42,3.14152,"fu":"bar"},"consent":true}}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPJSONDecoder()},
			expectedError: "The request was malformed or contained invalid parameters",
		},
		{
			d:       "should fail validation for wrong JSON type with unrestricted object",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(`{"dynamic_object":[42,3.14152]}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPJSONDecoder()},
			expectedError: "expected object, but got array",
		},
		{
			d:       "should accept JSON data with unrestricted object",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(`{"dynamic_object":{"stuff":{"blub":[42,3.14152],"fu":"bar"},"consent":true}}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPJSONDecoder()},
			expected: `{
	"dynamic_object": {
		"stuff": {
			"blub": [42, 3.14152],
			"fu": "bar"
		},
		"consent": true
	}
}`,
		},
		{
			d:       "should accept JSON data with unrestricted object and mixed object syntax and query parameter",
			request: newRequest(t, "POST", "/?name.last=Horstmann", bytes.NewBufferString(`{"dynamic_object":{"stuff":{"blub":[42,3.14152],"fu":"bar"},"consent":true},"name.first":"Horst"}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPJSONDecoder(),
				HTTPDecoderJSONFollowsFormFormat(),
				HTTPDecoderUseQueryAndBody()},
			expected: `{
	"dynamic_object": {
		"stuff": {
			"blub": [42, 3.14152],
			"fu": "bar"
		},
		"consent": true
	},
	"name": {
		"first": "Horst",
		"last": "Horstmann"
	}
}`,
		},
		{
			d:       "should accept JSON data with unrestricted object and mixed object syntax",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(`{"dynamic_object":{"stuff":{"blub":[42,3.14152],"fu":"bar"},"consent":true},"name.first":"Horst","name.last":"Horstmann"}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPJSONDecoder(),
				HTTPDecoderJSONFollowsFormFormat()},
			expected: `{
	"dynamic_object": {
		"stuff": {
			"blub": [42, 3.14152],
			"fu": "bar"
		},
		"consent": true
	},
	"name": {
		"first": "Horst",
		"last": "Horstmann"
	}
}`,
		},
		{
			d: "should fail form data with invalid premarshalled JSON object",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				"dynamic_object": {`{"stuff":{"blub":[42, 3.14152,"fu":"bar"},"consent":true}`},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPFormDecoder()},
			expectedError: "The request was malformed or contained invalid parameters",
		},
		{
			d: "should fail validation for form data with wrong premarshalled JSON type",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				"dynamic_object": {`[42, 3.14152]`},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPFormDecoder()},
			expectedError: "expected object, but got array",
		},
		{
			d: "should accept form data with premarshalled JSON object",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				"dynamic_object": {`{"stuff":{"blub":[42, 3.14152],"fu":"bar"},"consent":true}`},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPFormDecoder()},
			expected: `{
	"dynamic_object": {
		"stuff": {
			"blub": [42, 3.14152],
			"fu": "bar"
		},
		"consent": true
	},
	"name": {}
}`,
		},
		{
			d: "should accept form data with premarshalled JSON object and mixed object syntax and query parameter",
			request: newRequest(t, "POST", "/?name.last=Horstmann", bytes.NewBufferString(url.Values{
				"dynamic_object": {`{"stuff":{"blub":[42, 3.14152],"fu":"bar"},"consent":true}`},
				"name.first":     {"Horst"},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPFormDecoder(),
				HTTPDecoderUseQueryAndBody()},
			expected: `{
	"dynamic_object": {
		"stuff": {
			"blub": [42, 3.14152],
			"fu": "bar"
		},
		"consent": true
	},
	"name": {
		"first": "Horst",
		"last": "Horstmann"
	}
}`,
		},
		{
			d: "should accept form data with premarshalled JSON object and mixed object syntax",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				"dynamic_object": {`{"stuff":{"blub":[42, 3.14152],"fu":"bar"},"consent":true}`},
				"name.first":     {"Horst"},
				"name.last":      {"Horstmann"},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/dynamic-object.json", nil),
				HTTPFormDecoder()},
			expected: `{
	"dynamic_object": {
		"stuff": {
			"blub": [42, 3.14152],
			"fu": "bar"
		},
		"consent": true
	},
	"name": {
		"first": "Horst",
		"last": "Horstmann"
	}
}`,
		},
		{
			d: "should pass form request and type assert data",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				"name.first": {"Aeneas"},
				"name.last":  {"Rekkas"},
				"age":        {"29"},
				"ratio":      {"0.9"},
				"consent":    {"true"},

				// newsletter represents a special case for checkbox input with true/false and raw HTML.
				"newsletter": {
					"false", // comes from <input type="hidden" name="newsletter" value="false">
					"true",  // comes from <input type="checkbox" name="newsletter" value="true" checked>
				},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{HTTPJSONSchemaCompiler("stub/person.json", nil)},
			expected: `{
	"name": {"first": "Aeneas", "last": "Rekkas"},
	"age": 29,
	"newsletter": true,
	"consent": true,
	"ratio": 0.9
}`,
		},
		{
			d: "should mark the correct fields when nested objects are required",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				// newsletter represents a special case for checkbox input with true/false and raw HTML.
				"foo": {"bar"},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/consent.json", nil),
				HTTPKeepRequestBody(true),
				HTTPDecoderSetValidatePayloads(false),
				HTTPDecoderUseQueryAndBody(),
				HTTPDecoderAllowedMethods("POST", "GET"),
				HTTPDecoderJSONFollowsFormFormat(),
			},
			expected: `{
  "traits": {
	"consent": {
	  "inner": {}
    },
	"notrequired": {}
  }
}`,
		},
		{
			d: "should pass form request with payload in query and type assert data",
			request: newRequest(t, "POST", "/?age=29", bytes.NewBufferString(url.Values{
				"name.first": {"Aeneas"},
				"name.last":  {"Rekkas"},
				"ratio":      {"0.9"},
				"consent":    {"true"},
				// newsletter represents a special case for checkbox input with true/false and raw HTML.
				"newsletter": {
					"false", // comes from <input type="hidden" name="newsletter" value="false">
					"true",  // comes from <input type="checkbox" name="newsletter" value="true" checked>
				},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{HTTPJSONSchemaCompiler("stub/person.json", nil)},
			expected: `{
	"name": {"first": "Aeneas", "last": "Rekkas"},
	"newsletter": true,
	"consent": true,
	"ratio": 0.9
}`,
		},
		{
			d: "should pass form request with payload in query and type assert data",
			request: newRequest(t, "POST", "/?age=29", bytes.NewBufferString(url.Values{
				"name.first": {"Aeneas"},
				"name.last":  {"Rekkas"},
				"ratio":      {"0.9"},
				"consent":    {"true"},
				// newsletter represents a special case for checkbox input with true/false and raw HTML.
				"newsletter": {
					"false", // comes from <input type="hidden" name="newsletter" value="false">
					"true",  // comes from <input type="checkbox" name="newsletter" value="true" checked>
				},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPDecoderUseQueryAndBody(),
				HTTPJSONSchemaCompiler("stub/person.json", nil),
			},
			expected: `{
	"name": {"first": "Aeneas", "last": "Rekkas"},
	"age": 29,
	"newsletter": true,
	"consent": true,
	"ratio": 0.9
}`,
		},
		{
			d: "should fail form request if empty values are sent because of required fields",
			request: newRequest(t, "POST", "/?age=29", bytes.NewBufferString(url.Values{
				"name.first":  {""},
				"name.last":   {""},
				"name2.first": {""},
				"name2.last":  {""},
				"ratio":       {""},
				"ratio2":      {""},
				"age":         {""},
				"age2":        {""},
				"consent":     {""},
				"consent2":    {""},
				// newsletter represents a special case for checkbox input with true/false and raw HTML.
				"newsletter":  {""},
				"newsletter2": {""},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPDecoderUseQueryAndBody(),
				HTTPJSONSchemaCompiler("stub/required-defaults.json", nil),
			},
			expectedError: `I[#/name2] S[#/properties/name2/required] missing properties: "first"`,
		},
		{
			d:             "should fail json request formatted as form if payload is invalid",
			request:       newRequest(t, "POST", "/", bytes.NewBufferString(`{"name.first":"Aeneas", "name.last":"Rekkas","age":"not-a-number"}`), httpContentTypeJSON),
			options:       []HTTPDecoderOption{HTTPJSONSchemaCompiler("stub/person.json", nil)},
			expectedError: "expected integer, but got string",
		},
		{
			d: "should pass JSON request formatted as a form",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(`{
	"name.first": "Aeneas",
	"name.last":  "Rekkas",
	"age":        29,
	"ratio":      0.9,
	"consent":    false,
	"newsletter": true
}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{HTTPDecoderJSONFollowsFormFormat(),
				HTTPJSONSchemaCompiler("stub/person.json", nil)},
			expected: `{
	"name": {"first": "Aeneas", "last": "Rekkas"},
	"age": 29,
	"newsletter": true,
	"consent": false,
	"ratio": 0.9
}`,
		},
		{
			d: "should pass JSON request formatted as a form",
			request: newRequest(t, "POST", "/?age=29", bytes.NewBufferString(`{
	"name.first": "Aeneas",
	"name.last":  "Rekkas",
	"ratio":      0.9,
	"consent":    false,
	"newsletter": true
}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{HTTPDecoderJSONFollowsFormFormat(),
				HTTPJSONSchemaCompiler("stub/person.json", nil)},
			expected: `{
	"name": {"first": "Aeneas", "last": "Rekkas"},
	"newsletter": true,
	"consent": false,
	"ratio": 0.9
}`,
		},
		{
			d: "should pass JSON request formatted as a JSON even if HTTPDecoderJSONFollowsFormFormat is used",
			request: newRequest(t, "POST", "/?age=29", bytes.NewBufferString(`{
	"name": {"first": "Aeneas", "last": "Rekkas"},
	"ratio":      0.9,
	"consent":    false,
	"newsletter": true
}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{HTTPDecoderJSONFollowsFormFormat(),
				HTTPJSONSchemaCompiler("stub/person.json", nil)},
			expected: `{
	"name": {"first": "Aeneas", "last": "Rekkas"},
	"newsletter": true,
	"consent": false,
	"ratio": 0.9
}`,
		},
		{
			d: "should not retry indefinitely if key does not exist",
			request: newRequest(t, "POST", "/?age=29", bytes.NewBufferString(`{
	"not-foo": "bar"
}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{HTTPDecoderJSONFollowsFormFormat(),
				HTTPJSONSchemaCompiler("stub/schema.json", nil)},
			expectedError: "I[#] S[#/required] missing properties",
		},
		{
			d:       "should indicate the true missing fields from nested form",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{"leaf": {"foo"}}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPDecoderUseQueryAndBody(),
				HTTPDecoderSetIgnoreParseErrorsStrategy(ParseErrorIgnoreConversionErrors),
				HTTPJSONSchemaCompiler("stub/nested.json", nil)},
			expectedError: `I[#/node/node/node] S[#/properties/node/properties/node/properties/node/required] missing properties: "leaf"`,
		},
		{
			d: "should pass JSON request formatted as a form",
			request: newRequest(t, "POST", "/?age=29", bytes.NewBufferString(`{
	"name.first": "Aeneas",
	"name.last":  "Rekkas",
	"ratio":      0.9,
	"consent":    false,
	"newsletter": true
}`), httpContentTypeJSON),
			options: []HTTPDecoderOption{
				HTTPDecoderUseQueryAndBody(),
				HTTPDecoderJSONFollowsFormFormat(),
				HTTPJSONSchemaCompiler("stub/person.json", nil)},
			expected: `{
	"name": {"first": "Aeneas", "last": "Rekkas"},
	"age": 29,
	"newsletter": true,
	"consent": false,
	"ratio": 0.9
}`,
		},
		{
			d: "should pass JSON request GET request",
			request: newRequest(t, "GET", "/?"+url.Values{
				"name.first": {"Aeneas"},
				"name.last":  {"Rekkas"},
				"age":        {"29"},
				"ratio":      {"0.9"},
				"consent":    {"false"},
				"newsletter": {"true"},
			}.Encode(), nil, ""),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/person.json", nil),
				HTTPDecoderAllowedMethods("GET"),
			},
			expected: `{
	"name": {"first": "Aeneas", "last": "Rekkas"},
	"age": 29,
	"newsletter": true,
	"consent": false,
	"ratio": 0.9
}`,
		},
		{
			d:       "should fail because json is not an object when using form format",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(`[]`), httpContentTypeJSON),
			options: []HTTPDecoderOption{HTTPDecoderJSONFollowsFormFormat(),
				HTTPJSONSchemaCompiler("stub/person.json", nil)},
			expectedError: "be an object",
		},
		{
			d: "should work with ParseErrorIgnoreConversionErrors",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				"ratio": {"foobar"},
			}.Encode()), httpContentTypeURLEncodedForm),
			options: []HTTPDecoderOption{
				HTTPJSONSchemaCompiler("stub/person.json", nil),
				HTTPDecoderSetIgnoreParseErrorsStrategy(ParseErrorIgnoreConversionErrors),
				HTTPDecoderSetValidatePayloads(false),
			},
			expected: `{"name": {}, "ratio": "foobar"}`,
		},
		{
			d: "should work with ParseErrorIgnoreConversionErrors",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				"ratio": {"foobar"},
			}.Encode()), httpContentTypeURLEncodedForm),
			options:  []HTTPDecoderOption{HTTPJSONSchemaCompiler("stub/person.json", nil), HTTPDecoderSetIgnoreParseErrorsStrategy(ParseErrorUseEmptyValueOnConversionErrors)},
			expected: `{"name": {}, "ratio": 0.0}`,
		},
		{
			d: "should work with ParseErrorIgnoreConversionErrors",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				"ratio": {"foobar"},
			}.Encode()), httpContentTypeURLEncodedForm),
			options:       []HTTPDecoderOption{HTTPJSONSchemaCompiler("stub/person.json", nil), HTTPDecoderSetIgnoreParseErrorsStrategy(ParseErrorReturnOnConversionErrors)},
			expectedError: `strconv.ParseFloat: parsing "foobar"`,
		},
		{
			d: "should interpret numbers as string if mandated by the schema",
			request: newRequest(t, "POST", "/", bytes.NewBufferString(url.Values{
				"name.first": {"12345"},
			}.Encode()), httpContentTypeURLEncodedForm),
			options:  []HTTPDecoderOption{HTTPJSONSchemaCompiler("stub/person.json", nil), HTTPDecoderSetIgnoreParseErrorsStrategy(ParseErrorUseEmptyValueOnConversionErrors)},
			expected: `{"name": {"first": "12345"}}`,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			dec := NewHTTP()
			var destination json.RawMessage
			err := dec.Decode(tc.request, &destination, tc.options...)
			if tc.expectedError != "" {
				if e, ok := errors.Cause(err).(*jsonschema.ValidationError); ok {
					t.Logf("%+v", e)
				}
				require.Error(t, err)
				require.Contains(t, fmt.Sprintf("%+v", err), tc.expectedError)
				if len(tc.expected) > 0 {
					assert.JSONEq(t, tc.expected, string(destination))
				}
				return
			}

			require.NoError(t, err)
			assertx.EqualAsJSON(t, json.RawMessage(tc.expected), destination)
		})
	}

	t.Run("description=read body twice", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		dec := NewHTTP()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer wg.Done()

			var destination json.RawMessage
			require.NoError(t, dec.Decode(r, &destination, HTTPJSONSchemaCompiler("stub/person.json", nil), HTTPKeepRequestBody(true)))
			assert.EqualValues(t, "12345", gjson.GetBytes(destination, "name.first").String())

			require.NoError(t, dec.Decode(r, &destination, HTTPJSONSchemaCompiler("stub/person.json", nil), HTTPKeepRequestBody(true)))
			assert.EqualValues(t, "12345", gjson.GetBytes(destination, "name.first").String())
		}))
		t.Cleanup(ts.Close)

		_, err := ts.Client().PostForm(ts.URL, url.Values{"name.first": {"12345"}})
		require.NoError(t, err)

		wg.Wait()
	})
}
