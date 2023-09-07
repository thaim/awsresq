package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestEfsValidate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqEfsAPI
		resource string
		expected bool
	}{
		{
			name:     "validate file-system resource",
			api:      AwsresqEfsAPI{},
			resource: "file-system",
			expected: true,
		},
		{
			name:     "validate undefined resource",
			api:      AwsresqEfsAPI{},
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

func TestEfsFileSystemQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsEfsAPI(ctrl)

	mc.EXPECT().
		DescribeFileSystems(gomock.Any(), nil).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []types.FileSystemDescription{
				{
					FileSystemId: aws.String("fs-0123456789abcdef0"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.FileSystemDescription
		wantErr   bool
		expectErr string
	}{
		{
			name: "query cluster resource",
			expected: []types.FileSystemDescription{
				{
					FileSystemId: aws.String("fs-0123456789abcdef0"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqEfsAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("file-system")

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
				if !strings.Contains(err.Error(), tt.expectErr) {
					t.Errorf("expected %v, but got %v", tt.expectErr, err.Error())
				}
			}
			if err != nil {
				t.Errorf("expected nil, but got %v", err.Error())
			}

			if actual.Service != "efs" {
				t.Errorf("expected efs, but got %v", actual.Service)
			}
			if actual.Resource != "file-system" {
				t.Errorf("expected file-system, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.FileSystemDescription)
				if !ok {
					t.Errorf("expected types.FileSystemDescription, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.FileSystemId , tt.expected[i].FileSystemId) {
					t.Errorf("expected %v, but got %v", tt.expected[i].FileSystemId, actualOutput.FileSystemId)
				}
			}
		})
	}
}
