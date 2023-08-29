package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestLogsValidate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqLogsAPI
		resource string
		expected bool
	}{
		{
			name:     "validate log-group resource",
			api:      AwsresqLogsAPI{},
			resource: "log-group",
			expected: true,
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

func TestLogsQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsLogsAPI(ctrl)

	mc.EXPECT().
		DescribeLogGroups(gomock.Any(), nil).
		Return(&cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []types.LogGroup{
				{
					LogGroupName: aws.String("/aws/lambda/test-lambda01"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		resource  string
		expected  *cloudwatchlogs.DescribeLogGroupsOutput
		wantErr   bool
		expectErr string
	}{
		{
			name:     "query log-group resource",
			resource: "log-group",
			expected: &cloudwatchlogs.DescribeLogGroupsOutput{
				LogGroups: []types.LogGroup{
					{
						LogGroupName: aws.String("/aws/lambda/test-lambda01"),
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadDefaultConfig(context.Background())
			if err != nil {
				t.Errorf("failed to load config: %v", err)
			}
			api := NewAwsresqLogsAPI(cfg, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query(tt.resource)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
				if err.Error() != tt.expectErr {
					t.Errorf("expected %v, but got %v", tt.expectErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}

			if actual.Service != "logs" {
				t.Errorf("expected logs, but got %v", actual.Service)
			}
			if actual.Resource != tt.resource {
				t.Errorf("expected %v, but got %v", tt.resource, actual.Resource)
			}
			if len(tt.expected.LogGroups) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected.LogGroups), len(actual.Results))
			}
			for i := range tt.expected.LogGroups {
				actualOutput, ok := actual.Results[i].(types.LogGroup)
				if !ok {
					t.Errorf("expected types.LogGroup, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(tt.expected.LogGroups[i], actualOutput) {
					t.Errorf("expected %+v, but got %+v", tt.expected.LogGroups[i], actualOutput)
				}
			}
		})
	}
}
