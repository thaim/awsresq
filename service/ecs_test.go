package service

import (
	"testing"
)


func TestEcsValidate(t *testing.T) {
	cases := []struct {
		name string
		api AwsresqEcsAPI
		resource string
		expected bool
	}{
		{
			name: "validate function resource",
			api: AwsresqEcsAPI{},
			resource: "task-definition",
			expected: true,
		},
		{
			name: "validate undefined resource",
			api: AwsresqEcsAPI{},
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
