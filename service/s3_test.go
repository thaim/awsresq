package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestS3Validate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqS3API
		resource string
		expected bool
	}{
		{
			name:     "validate bucket resource",
			api:      AwsresqS3API{},
			resource: "bucket",
			expected: true,
		},
		{
			name:     "validate undefined resource",
			api:      AwsresqS3API{},
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

func TestS3BucketQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsS3API(ctrl)

	mc.EXPECT().
		ListBuckets(gomock.Any(), nil).
		Return(&s3.ListBucketsOutput{
			Buckets: []types.Bucket{
				{
					Name: aws.String("test-bucket"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.Bucket
		wantErr   bool
		expectErr string
	}{
		{
			name: "query cluster resource",
			expected: []types.Bucket{
				{
					Name: aws.String("test-bucket"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqS3API(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("bucket")

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

			if actual.Service != "s3" {
				t.Errorf("expected s3, but got %v", actual.Service)
			}
			if actual.Resource != "bucket" {
				t.Errorf("expected bucket, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Bucket)
				if !ok {
					t.Errorf("expected types.FileSystemDescription, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.Name, tt.expected[i].Name) {
					t.Errorf("expected %v, but got %v", *tt.expected[i].Name, *actualOutput.Name)
				}
			}
		})
	}
}
