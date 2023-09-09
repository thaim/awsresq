package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestEcrValidate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqEcrAPI
		resource string
		expected bool
	}{
		{
			name:     "validate repository resource",
			api:      AwsresqEcrAPI{},
			resource: "repository",
			expected: true,
		},
		{
			name:     "validate undefined resource",
			api:      AwsresqEcrAPI{},
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

func TestEcrRepositoryQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsEcrAPI(ctrl)

	mc.EXPECT().
		DescribeRepositories(gomock.Any(), nil).
		Return(&ecr.DescribeRepositoriesOutput{
			Repositories: []types.Repository{
				{
					RepositoryName: aws.String("test"),
				},
			},
		}, nil)

	cases := []struct {
		name     string
		expected []types.Repository
		wantErr  bool
		expectErr string
	}{
		{
			name: "query repository",
			expected: []types.Repository{
				{
					RepositoryName: aws.String("test"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.Background())
			api := NewAwsresqEcrAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("repository")

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

			if actual.Service != "ecr" {
				t.Errorf("expected ecr, but got %v", actual.Service)
			}
			if actual.Resource != "repository" {
				t.Errorf("expected repository, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Repository)
				if !ok {
					t.Errorf("expected types.Repository, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.RepositoryName, tt.expected[i].RepositoryName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].RepositoryName, actualOutput.RepositoryName)
				}
			}
		})
	}
}
