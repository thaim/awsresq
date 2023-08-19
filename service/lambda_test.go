package service

import (
	"testing"
)


func TestValidate(t *testing.T) {
	cases := []struct {
		name string
		api AwsLambdaAPI
		resource string
		expected bool
	}{
		{
			name: "validate function resource",
			api: AwsLambdaAPI{},
			resource: "function",
			expected: true,
		},
		{
			name: "validate undefined resource",
			api: AwsLambdaAPI{},
			resource: "undefined",
			expected: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.api.Validate(tt.resource)

			if actual != tt.expected {
				t.Errorf("expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}
