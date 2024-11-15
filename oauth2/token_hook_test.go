package oauth2

import (
	"reflect"
	"testing"
)

func TestUpdateExtraClaims(t *testing.T) {
	tests := []struct {
		name               string
		priorExtraClaims   map[string]interface{}
		webhookExtraClaims map[string]interface{}
		expected           map[string]interface{}
	}{
		{
			name: "Merge with no updates",
			priorExtraClaims: map[string]interface{}{
				"claim1": "value1",
				"claim2": "value2",
			},
			webhookExtraClaims: map[string]interface{}{
				"claim3": "value3",
				"claim4": "value4",
			},
			expected: map[string]interface{}{
				"claim1": "value1",
				"claim2": "value2",
				"claim3": "value3",
				"claim4": "value4",
			},
		},
		{
			name: "Merge with updates",
			priorExtraClaims: map[string]interface{}{
				"claim1": "value1",
				"claim2": "value2",
			},
			webhookExtraClaims: map[string]interface{}{
				"claim2": "newValue2", // Overwrites prior claim2
				"claim3": "value3",
			},
			expected: map[string]interface{}{
				"claim1": "value1",
				"claim2": "newValue2",
				"claim3": "value3",
			},
		},
		{
			name: "Empty webhook claims",
			priorExtraClaims: map[string]interface{}{
				"claim1": "value1",
			},
			webhookExtraClaims: map[string]interface{}{},
			expected: map[string]interface{}{
				"claim1": "value1",
			},
		},
		{
			name:             "Empty prior claims",
			priorExtraClaims: map[string]interface{}{},
			webhookExtraClaims: map[string]interface{}{
				"claim1": "value1",
			},
			expected: map[string]interface{}{
				"claim1": "value1",
			},
		},
		{
			name:               "Both maps empty",
			priorExtraClaims:   map[string]interface{}{},
			webhookExtraClaims: map[string]interface{}{},
			expected:           map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := updateExtraClaims(tt.priorExtraClaims, tt.webhookExtraClaims)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("updateExtraClaims() = %v, want %v", result, tt.expected)
			}
		})
	}
}
