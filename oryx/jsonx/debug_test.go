// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/x/jsonx"
)

func TestJSONShape(t *testing.T) {
	for _, tc := range []struct {
		name     string
		in       string
		expected string
	}{{
		name: "user patch",
		in: `{
  "schemas" : [ "urn:ietf:params:scim:schemas:core:2.0:User" ],
  "id" : "d4b4f9db-2361-4845-a4cd-51e12527b92e",
  "externalId" : "00uo3xq5f75s2KCOE5d7",
  "userName" : "henning.perl@ory.sh",
  "name" : {
    "familyName" : "Perl",
    "givenName" : "Henning"
  },
  "displayName" : "Henning Perl",
  "locale" : "en-US",
  "active" : true,
  "emails" : [ {
    "value" : "henning.perl@ory.sh",
    "primary" : true,
    "type" : "work"
  } ],
  "groups" : [ {
    "value" : "21c5f2f9-8fb0-45b3-9bb6-61ecd1090549",
    "display" : "Developers",
    "type" : "direct"
  }, {
    "value" : "a37d499d-739c-4e08-8273-c124f85172fe",
    "display" : "SCIM pros",
    "type" : "direct"
  } ],
  "meta" : {
    "resourceType" : "User",
    "created" : "2025-04-25T07:53:43Z",
    "lastModified" : "2025-04-25T08:31:23Z"
  },
  "roles" : [ "foo", "bar" ]
}`,
		expected: `{
  "active": "boolean",
  "displayName": "string",
  "emails": [
	{
	  "primary": "boolean",
	  "type": "string",
	  "value": "string"
	}
  ],
  "externalId": "string",
  "groups": [
	{
	  "display": "string",
	  "type": "string",
	  "value": "string"
	},
	{
	  "display": "string",
	  "type": "string",
	  "value": "string"
	}
  ],
  "id": "d4b4f9db-2361-4845-a4cd-51e12527b92e",
  "locale": "string",
  "meta": {
	"created": "string",
	"lastModified": "string",
	"resourceType": "string"
  },
  "name": {
	"familyName": "string",
	"givenName": "string"
  },
  "roles": [
	"string",
	"string"
  ],
  "schemas": [
	"urn:ietf:params:scim:schemas:core:2.0:User"
  ],
  "userName": "string"
}`,
	}, {
		name:     "invalid JSON",
		in:       `{`,
		expected: `{"error": "invalid JSON", "message": "unexpected end of JSON input"}`,
	}, {
		name: "different types",
		in: `{
	"float": 0.42,
	"int": 42,
	"string": "foo",
	"bool": true,
	"null": null,
	"array": [1, "2", 0]
}`,
		expected: `{
	"float": "number",
	"int": "number",
	"string": "string",
	"bool": "boolean",
	"null": "null",
	"array": ["number", "string", "number"]
}`,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			actual := string(jsonx.Anonymize([]byte(tc.in), "id", "schemas"))
			assert.JSONEq(t, tc.expected, actual, actual)
		})
	}
}
