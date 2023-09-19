package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestCloudWatchValidate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqCloudwatchAPI
		resource string
		expect   bool
	}{
		{
			name:     "valid metric resource",
			api:      AwsresqCloudwatchAPI{},
			resource: "metric",
			expect:   true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.api.Validate(tt.resource)

			if actual != tt.expect {
				t.Errorf("expected %t, got %t", tt.expect, actual)
			}
		})
	}
}

func TestCloudwatchMetricQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsCloudwatchAPI(ctrl)

	mc.EXPECT().
		ListMetrics(gomock.Any(), nil).
		Return(&cloudwatch.ListMetricsOutput{
			Metrics: []types.Metric{
				{
					MetricName: aws.String("CPUUtilization"),
					Namespace:  aws.String("AWS/ECS"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.Metric
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid metric resource",
			expected: []types.Metric{
				{
					MetricName: aws.String("CPUUtilization"),
					Namespace:  aws.String("AWS/ECS"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqCloudwatchAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("metric")

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

			if actual.Service != "cloudwatch" {
				t.Errorf("expected cloudwatch, but got %v", actual.Service)
			}
			if actual.Resource != "metric" {
				t.Errorf("expected metric, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Metric)
				if !ok {
					t.Errorf("expected types.Metric, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.MetricName, tt.expected[i].MetricName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].MetricName, actualOutput.MetricName)
				}
				if !reflect.DeepEqual(actualOutput.Namespace, tt.expected[i].Namespace) {
					t.Errorf("expected %v, but got %v", tt.expected[i].Namespace, actualOutput.Namespace)
				}
			}
		})
	}
}
