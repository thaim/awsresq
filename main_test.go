package main

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	cases := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "value",
			version:  "0.1.0",
			expected: "0.1.0",
		},
		{
			name:     "value",
			version:  "",
			expected: "(devel)",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			defaultVersion := version
			version = tt.version
			actual := getVersion()

			if actual != tt.expected {
				t.Errorf("getVersion() = %v, want %v", actual, tt.expected)
			}
			version = defaultVersion
		})
	}
}
