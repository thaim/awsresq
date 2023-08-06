package internal

import (
	"reflect"
	"testing"
)

func TestBuildRegion(t *testing.T) {
	cases := []struct {
		name string
		input string
		expected []string
	}{
		{
			name: "build regions from all",
			input: "all",
			expected: []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ap-south-1", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "cn-north-1", "cn-nothwest-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-west-3", "sa-east-1"},
		},
		{
			name: "specify single region",
			input: "ap-northeast-1",
			expected: []string{"ap-northeast-1"},
		},
		{
			name: "specify multiple regions",
			input: "ap-northeast-1,us-east-1,us-west-1",
			expected: []string{"ap-northeast-1", "us-east-1", "us-west-1"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual := buildRegion(tt.input)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("actual = %v, want = %v", actual, tt.expected)
			}
		})
	}
}
