// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package decoderx

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/ory/jsonschema/v3"

	"github.com/ory/herodot"

	"github.com/ory/x/httpx"
	"github.com/ory/x/jsonschemax"
)

type (
	// HTTP decodes json and form-data from HTTP Request Bodies.
	HTTP struct{}

	httpDecoderOptions struct {
		keepRequestBody           bool
		allowedContentTypes       []string
		allowedHTTPMethods        []string
		jsonSchemaRef             string
		jsonSchemaCompiler        *jsonschema.Compiler
		jsonSchemaValidate        bool
		maxCircularReferenceDepth uint8
		handleParseErrors         parseErrorStrategy
		expectJSONFlattened       bool
		queryAndBody              bool
	}

	// HTTPDecoderOption configures the HTTP decoder.
	HTTPDecoderOption func(*httpDecoderOptions)

	parseErrorStrategy uint8
)

const (
	httpContentTypeMultipartForm  = "multipart/form-data"
	httpContentTypeURLEncodedForm = "application/x-www-form-urlencoded"
	httpContentTypeJSON           = "application/json"
)

const (
	// ParseErrorIgnoreConversionErrors will ignore any errors caused by strconv.Parse* and use the
	// raw form field value, which is a string, when such a parse error occurs.
	//
	// If the JSON Schema defines `{"ratio": {"type": "number"}}` but `ratio=foobar` then field
	// `ratio` will be handled as a string. If the destination struct is a `json.RawMessage`, then
	// the output will be `{"ratio": "foobar"}`.
	ParseErrorIgnoreConversionErrors parseErrorStrategy = iota + 1

	// ParseErrorUseEmptyValueOnConversionErrors will ignore any parse errors caused by strconv.Parse* and use the
	// default value of the type to be casted, e.g. float64(0), string("").
	//
	// If the JSON Schema defines `{"ratio": {"type": "number"}}` but `ratio=foobar` then field
	// `ratio` will receive the default value for the primitive type (here `0.0` for `number`).
	// If the destination struct is a `json.RawMessage`, then the output will be `{"ratio": 0.0}`.
	ParseErrorUseEmptyValueOnConversionErrors

	// ParseErrorReturnOnConversionErrors will abort and return with an error if strconv.Parse* returns
	// an error.
	//
	// If the JSON Schema defines `{"ratio": {"type": "number"}}` but `ratio=foobar` the parser aborts
	// and returns an error, here: `strconv.ParseFloat: parsing "foobar"`.
	ParseErrorReturnOnConversionErrors
)

var errKeyNotFound = errors.New("key not found")

// HTTPFormDecoder configures the HTTP decoder to only accept form-data
// (application/x-www-form-urlencoded, multipart/form-data)
func HTTPFormDecoder() HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		o.allowedContentTypes = []string{httpContentTypeMultipartForm, httpContentTypeURLEncodedForm}
	}
}

// HTTPJSONDecoder configures the HTTP decoder to only accept JSON data
// (application/json).
func HTTPJSONDecoder() HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		o.allowedContentTypes = []string{httpContentTypeJSON}
	}
}

// HTTPKeepRequestBody configures the HTTP decoder to allow other
// HTTP request body readers to read the body as well by keeping
// the data in memory.
func HTTPKeepRequestBody(keep bool) HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		o.keepRequestBody = keep
	}
}

// HTTPDecoderSetValidatePayloads sets if payloads should be validated or not.
func HTTPDecoderSetValidatePayloads(validate bool) HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		o.jsonSchemaValidate = validate
		o.keepRequestBody = true
	}
}

// HTTPDecoderJSONFollowsFormFormat if set tells the decoder that JSON follows the same conventions
// as the form decoder, meaning `{"foo.bar": "..."}` is translated to `{"foo": {"bar": "..."}}`.
func HTTPDecoderJSONFollowsFormFormat() HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		o.expectJSONFlattened = true
		o.keepRequestBody = true
	}
}

// HTTPDecoderAllowedMethods sets the allowed HTTP methods. Defaults are POST, PUT, PATCH.
func HTTPDecoderAllowedMethods(method ...string) HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		o.allowedHTTPMethods = method
	}
}

// HTTPDecoderUseQueryAndBody will check both the HTTP body and the HTTP query params when decoding.
// Only relevant for non-GET operations.
func HTTPDecoderUseQueryAndBody() HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		o.queryAndBody = true
	}
}

// HTTPDecoderSetIgnoreParseErrorsStrategy sets a strategy for dealing with strconv.Parse* errors:
//
// - decoderx.ParseErrorIgnoreConversionErrors will ignore any parse errors caused by strconv.Parse* and use the
// raw form field value, which is a string, when such a parse error occurs. (default)
// - decoderx.ParseErrorUseEmptyValueOnConversionErrors will ignore any parse errors caused by strconv.Parse* and use the
// default value of the type to be casted, e.g. float64(0), string("").
// - decoderx.ParseErrorReturnOnConversionErrors will abort and return with an error if strconv.Parse* returns
// an error.
func HTTPDecoderSetIgnoreParseErrorsStrategy(strategy parseErrorStrategy) HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		o.handleParseErrors = strategy
	}
}

// HTTPDecoderSetMaxCircularReferenceDepth sets the maximum recursive reference resolution depth.
func HTTPDecoderSetMaxCircularReferenceDepth(depth uint8) HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		o.maxCircularReferenceDepth = depth
	}
}

// HTTPJSONSchemaCompiler sets a JSON schema to be used for validation and type assertion of
// incoming requests.
func HTTPJSONSchemaCompiler(ref string, compiler *jsonschema.Compiler) HTTPDecoderOption {
	return func(o *httpDecoderOptions) {
		if compiler == nil {
			compiler = jsonschema.NewCompiler()
		}
		compiler.ExtractAnnotations = true
		o.jsonSchemaCompiler = compiler
		o.jsonSchemaRef = ref
		o.jsonSchemaValidate = true
	}
}

// HTTPRawJSONSchemaCompiler uses a JSON Schema Compiler with the provided JSON Schema in raw byte form.
func HTTPRawJSONSchemaCompiler(raw []byte) (HTTPDecoderOption, error) {
	compiler := jsonschema.NewCompiler()
	id := fmt.Sprintf("%x.json", sha256.Sum256(raw))
	if err := compiler.AddResource(id, bytes.NewReader(raw)); err != nil {
		return nil, err
	}
	compiler.ExtractAnnotations = true

	return func(o *httpDecoderOptions) {
		o.jsonSchemaCompiler = compiler
		o.jsonSchemaRef = id
		o.jsonSchemaValidate = true
	}, nil
}

// MustHTTPRawJSONSchemaCompiler uses HTTPRawJSONSchemaCompiler and panics on error.
func MustHTTPRawJSONSchemaCompiler(raw []byte) HTTPDecoderOption {
	f, err := HTTPRawJSONSchemaCompiler(raw)
	if err != nil {
		panic(err)
	}
	return f
}

func newHTTPDecoderOptions(fs []HTTPDecoderOption) *httpDecoderOptions {
	o := &httpDecoderOptions{
		allowedContentTypes: []string{
			httpContentTypeMultipartForm, httpContentTypeURLEncodedForm, httpContentTypeJSON,
		},
		allowedHTTPMethods:        []string{"POST", "PUT", "PATCH"},
		maxCircularReferenceDepth: 5,
		handleParseErrors:         ParseErrorIgnoreConversionErrors,
	}

	for _, f := range fs {
		f(o)
	}

	return o
}

// NewHTTP creates a new HTTP decoder.
func NewHTTP() *HTTP {
	return new(HTTP)
}

func (t *HTTP) validateRequest(r *http.Request, c *httpDecoderOptions) error {
	method := strings.ToUpper(r.Method)

	if !slices.Contains(c.allowedHTTPMethods, method) {
		return errors.WithStack(herodot.ErrBadRequest.WithReasonf(`Unable to decode body because HTTP Request Method was "%s" but only %v are supported.`, method, c.allowedHTTPMethods))
	}

	if method != "GET" {
		if r.ContentLength == 0 {
			return errors.WithStack(herodot.ErrBadRequest.WithReasonf(`Unable to decode HTTP Request Body because its HTTP Header "Content-Length" is zero.`))
		}

		if !httpx.HasContentType(r, c.allowedContentTypes...) {
			return errors.WithStack(herodot.ErrBadRequest.WithReasonf(`HTTP %s Request used unknown HTTP Header "Content-Type: %s", only %v are supported.`, method, r.Header.Get("Content-Type"), c.allowedContentTypes))
		}
	}

	return nil
}

func (t *HTTP) validatePayload(ctx context.Context, raw json.RawMessage, c *httpDecoderOptions) error {
	if !c.jsonSchemaValidate {
		return nil
	}

	if c.jsonSchemaCompiler == nil {
		return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("JSON Schema Validation is required but no compiler was provided."))
	}

	schema, err := c.jsonSchemaCompiler.Compile(ctx, c.jsonSchemaRef)
	if err != nil {
		return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to load JSON Schema from location: %s", c.jsonSchemaRef).WithDebug(err.Error()))
	}

	if err := schema.Validate(bytes.NewBuffer(raw)); err != nil {
		if _, ok := err.(*jsonschema.ValidationError); ok {
			return errors.WithStack(err)
		}
		return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to process JSON Schema and input: %s", err).WithDebug(err.Error()))
	}

	return nil
}

// Decode takes a HTTP Request Body and decodes it into destination.
func (t *HTTP) Decode(r *http.Request, destination interface{}, opts ...HTTPDecoderOption) error {
	c := newHTTPDecoderOptions(opts)
	if err := t.validateRequest(r, c); err != nil {
		return err
	}

	if r.Method == "GET" {
		return t.decodeForm(r, destination, c)
	} else if httpx.HasContentType(r, httpContentTypeJSON) {
		if c.expectJSONFlattened {
			return t.decodeJSONForm(r, destination, c)
		}
		return t.decodeJSON(r, destination, c)
	} else if httpx.HasContentType(r, httpContentTypeMultipartForm, httpContentTypeURLEncodedForm) {
		return t.decodeForm(r, destination, c)
	}

	return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to determine decoder for content type: %s", r.Header.Get("Content-Type")))
}

func (t *HTTP) requestBody(r *http.Request, o *httpDecoderOptions) (reader io.ReadCloser, err error) {
	if strings.ToUpper(r.Method) == "GET" {
		return io.NopCloser(bytes.NewBufferString(r.URL.Query().Encode())), nil
	}

	if !o.keepRequestBody {
		return r.Body, nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read body")
	}

	_ = r.Body.Close() //  must close
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return io.NopCloser(bytes.NewBuffer(bodyBytes)), nil
}

func (t *HTTP) decodeJSONForm(r *http.Request, destination interface{}, o *httpDecoderOptions) error {
	if o.jsonSchemaCompiler == nil {
		return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to decode HTTP Form Body because no validation schema was provided. This is a code bug."))
	}

	paths, err := jsonschemax.ListPathsWithRecursion(r.Context(), o.jsonSchemaRef, o.jsonSchemaCompiler, o.maxCircularReferenceDepth)
	if err != nil {
		return errors.WithStack(herodot.ErrInternalServerError.WithTrace(err).WithReasonf("Unable to prepare JSON Schema for HTTP Post Body Form parsing: %s", err).WithDebugf("%+v", err))
	}

	reader, err := t.requestBody(r, o)
	if err != nil {
		return err
	}

	var interim json.RawMessage
	if err := json.NewDecoder(reader).Decode(&interim); err != nil {
		return errors.WithStack(herodot.ErrBadRequest.WithError(err.Error()).WithReason("Unable to decode form as JSON."))
	}

	parsed := gjson.ParseBytes(interim)
	if !parsed.IsObject() {
		return errors.WithStack(herodot.ErrBadRequest.WithReasonf("Expected JSON sent in request body to be an object but got: %s", parsed.Type.String()))
	}

	values := url.Values{}
	parsed.ForEach(func(k, v gjson.Result) bool {
		values.Set(k.String(), v.String())
		return true
	})

	if o.queryAndBody {
		_ = r.ParseForm()
		for k := range r.Form {
			values.Set(k, r.Form.Get(k))
		}
	}

	raw, err := t.decodeURLValues(values, paths, o)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, destination); err != nil {
		return errors.WithStack(err)
	}

	return t.validatePayload(r.Context(), raw, o)
}

func (t *HTTP) decodeForm(r *http.Request, destination interface{}, o *httpDecoderOptions) error {
	if o.jsonSchemaCompiler == nil {
		return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to decode HTTP Form Body because no validation schema was provided. This is a code bug."))
	}

	reader, err := t.requestBody(r, o)
	if err != nil {
		return err
	}

	defer func() {
		r.Body = reader
	}()

	if err := r.ParseForm(); err != nil {
		return errors.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to decode HTTP %s form body: %s", strings.ToUpper(r.Method), err).WithDebug(err.Error()))
	}

	paths, err := jsonschemax.ListPathsWithRecursion(r.Context(), o.jsonSchemaRef, o.jsonSchemaCompiler, o.maxCircularReferenceDepth)
	if err != nil {
		return errors.WithStack(herodot.ErrInternalServerError.WithTrace(err).WithReasonf("Unable to prepare JSON Schema for HTTP Post Body Form parsing: %s", err).WithDebugf("%+v", err))
	}

	values := r.PostForm
	if r.Method == "GET" || o.queryAndBody {
		values = r.Form
	}

	raw, err := t.decodeURLValues(values, paths, o)
	if err != nil && !errors.Is(err, errKeyNotFound) {
		return err
	}

	if err := json.NewDecoder(bytes.NewReader(raw)).Decode(destination); err != nil {
		return errors.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to decode JSON payload: %s", err))
	}

	return t.validatePayload(r.Context(), raw, o)
}

func (t *HTTP) decodeURLValues(values url.Values, paths []jsonschemax.Path, o *httpDecoderOptions) (json.RawMessage, error) {
	raw := json.RawMessage(`{}`)
	for key := range values {
		for _, path := range paths {
			if key == path.Name {
				var err error
				switch path.Type.(type) {
				case []string:
					raw, err = sjson.SetBytes(raw, path.Name, values[key])
				case []float64:
					for k, v := range values[key] {
						var f float64
						if f, err = strconv.ParseFloat(v, 64); err != nil {
							switch o.handleParseErrors {
							case ParseErrorIgnoreConversionErrors:
								raw, err = sjson.SetBytes(raw, path.Name+"."+strconv.Itoa(k), v)
							case ParseErrorUseEmptyValueOnConversionErrors:
								raw, err = sjson.SetBytes(raw, path.Name+"."+strconv.Itoa(k), f)
							case ParseErrorReturnOnConversionErrors:
								return nil, errors.WithStack(herodot.ErrBadRequest.WithReasonf("Expected value to be a number.").
									WithDetail("parse_error", err.Error()).
									WithDetail("name", key).
									WithDetailf("index", "%d", k).
									WithDetail("value", v))
							}
						} else {
							raw, err = sjson.SetBytes(raw, path.Name+"."+strconv.Itoa(k), f)
						}
					}
				case []bool:
					for k, v := range values[key] {
						var b bool
						if b, err = strconv.ParseBool(v); err != nil {
							switch o.handleParseErrors {
							case ParseErrorIgnoreConversionErrors:
								raw, err = sjson.SetBytes(raw, path.Name+"."+strconv.Itoa(k), v)
							case ParseErrorUseEmptyValueOnConversionErrors:
								raw, err = sjson.SetBytes(raw, path.Name+"."+strconv.Itoa(k), b)
							case ParseErrorReturnOnConversionErrors:
								return nil, errors.WithStack(herodot.ErrBadRequest.WithReasonf("Expected value to be a boolean.").
									WithDetail("parse_error", err.Error()).
									WithDetail("name", key).
									WithDetailf("index", "%d", k).
									WithDetail("value", v))
							}
						} else {
							raw, err = sjson.SetBytes(raw, path.Name+"."+strconv.Itoa(k), b)
						}
					}
				case []interface{}:
					raw, err = sjson.SetBytes(raw, path.Name, values[key])
				case bool:
					v := values[key][len(values[key])-1]
					if len(v) == 0 {
						if !path.Required {
							continue
						}
						v = "false"
					}

					var b bool
					if b, err = strconv.ParseBool(v); err != nil {
						switch o.handleParseErrors {
						case ParseErrorIgnoreConversionErrors:
							raw, err = sjson.SetBytes(raw, path.Name, v)
						case ParseErrorUseEmptyValueOnConversionErrors:
							raw, err = sjson.SetBytes(raw, path.Name, b)
						case ParseErrorReturnOnConversionErrors:
							return nil, errors.WithStack(herodot.ErrBadRequest.WithReasonf("Expected value to be a boolean.").
								WithDetail("parse_error", err.Error()).
								WithDetail("name", key).
								WithDetail("value", values.Get(key)))
						}
					} else {
						raw, err = sjson.SetBytes(raw, path.Name, b)
					}
				case float64:
					v := values.Get(key)
					if len(v) == 0 {
						if !path.Required {
							continue
						}
						v = "0.0"
					}

					var f float64
					if f, err = strconv.ParseFloat(v, 64); err != nil {
						switch o.handleParseErrors {
						case ParseErrorIgnoreConversionErrors:
							raw, err = sjson.SetBytes(raw, path.Name, v)
						case ParseErrorUseEmptyValueOnConversionErrors:
							raw, err = sjson.SetBytes(raw, path.Name, f)
						case ParseErrorReturnOnConversionErrors:
							return nil, errors.WithStack(herodot.ErrBadRequest.WithReasonf("Expected value to be a number.").
								WithDetail("parse_error", err.Error()).
								WithDetail("name", key).
								WithDetail("value", values.Get(key)))
						}
					} else {
						raw, err = sjson.SetBytes(raw, path.Name, f)
					}
				case string:
					v := values.Get(key)
					if len(v) == 0 {
						continue
					}

					raw, err = sjson.SetBytes(raw, path.Name, v)
				case map[string]interface{}:
					v := values.Get(key)
					if len(v) == 0 && !path.Required {
						continue
					}

					raw, err = sjson.SetRawBytes(raw, path.Name, []byte(v))
				case []map[string]interface{}:
					raw, err = sjson.SetBytes(raw, path.Name, values[key])
				}

				if err != nil {
					return nil, errors.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to type assert values from HTTP Post Body: %s", err))
				}
				break
			}
		}
	}

	for _, path := range paths {
		if path.TypeHint != jsonschemax.JSON {
			continue
		}

		if !gjson.GetBytes(raw, path.Name).Exists() {
			var err error
			raw, err = sjson.SetRawBytes(raw, path.Name, []byte(`{}`))
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}

	return raw, nil
}

func (t *HTTP) decodeJSON(r *http.Request, destination interface{}, o *httpDecoderOptions) error {
	reader, err := t.requestBody(r, o)
	if err != nil {
		return err
	}

	raw, err := io.ReadAll(reader)
	if err != nil {
		return errors.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to read HTTP POST body: %s", err))
	}

	dc := json.NewDecoder(bytes.NewReader(raw))
	if err := dc.Decode(destination); err != nil {
		return errors.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to decode JSON payload: %s", err).WithDebugf("Received request body: %s", string(raw)))
	}

	if err := t.validatePayload(r.Context(), raw, o); err != nil {
		if o.expectJSONFlattened && strings.Contains(err.Error(), "json: unknown field") {
			return t.decodeJSONForm(r, destination, o)
		}
		return err
	}

	return nil
}
