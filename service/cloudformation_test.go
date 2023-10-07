package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestCloudformationValidate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqCloudformationAPI
		resource string
		expected bool
	}{
		{
			name:     "validate stack resource",
			api:      AwsresqCloudformationAPI{},
			resource: "stack",
			expected: true,
		},
		{
			name:     "validate stack-set resource",
			api:      AwsresqCloudformationAPI{},
			resource: "stack-set",
			expected: true,
		},
		{
			name:     "validate undefined resource",
			api:      AwsresqCloudformationAPI{},
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

func TestCloudformationStackQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsCloudformationAPI(ctrl)

	mc.EXPECT().
		DescribeStacks(gomock.Any(), nil).
		Return(&cloudformation.DescribeStacksOutput{
			Stacks: []types.Stack{
				{
					StackId:     aws.String("arn:aws:cloudformation:us-east-1:123456789012:stack/stack-name/guid"),
					StackName:   aws.String("stack-name"),
					StackStatus: types.StackStatusCreateComplete,
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.Stack
		wantErr   bool
		expectErr string
	}{
		{
			name: "query stack resource",
			expected: []types.Stack{
				{
					StackId:     aws.String("arn:aws:cloudformation:us-east-1:123456789012:stack/stack-name/guid"),
					StackName:   aws.String("stack-name"),
					StackStatus: types.StackStatusCreateComplete,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqCloudformationAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("stack")

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

			if actual.Service != "cloudformation" {
				t.Errorf("expected cloudformation, but got %v", actual.Service)
			}
			if actual.Resource != "stack" {
				t.Errorf("expected stack, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Stack)
				if !ok {
					t.Errorf("expected types.Stack, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.StackId, tt.expected[i].StackId) {
					t.Errorf("expected %v, but got %v", tt.expected[i].StackId, actualOutput.StackId)
				}
				if !reflect.DeepEqual(actualOutput.StackName, tt.expected[i].StackName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].StackName, actualOutput.StackName)
				}
				if !reflect.DeepEqual(actualOutput.StackStatus, tt.expected[i].StackStatus) {
					t.Errorf("expected %v, but got %v", tt.expected[i].StackStatus, actualOutput.StackStatus)
				}
			}
		})
	}
}

func TestCloudformationStackSetQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsCloudformationAPI(ctrl)

	mc.EXPECT().
		ListStackSets(gomock.Any(), nil).
		Return(&cloudformation.ListStackSetsOutput{
			Summaries: []types.StackSetSummary{
				{
					StackSetId:   aws.String("arn:aws:cloudformation:us-east-1:123456789012:stackset/stack-set-name/guid"),
					StackSetName: aws.String("stack-set-name"),
				},
			},
		}, nil).
		AnyTimes()

	mc.EXPECT().
		DescribeStackSet(gomock.Any(), &cloudformation.DescribeStackSetInput{
			StackSetName: aws.String("stack-set-name"),
		}).
		Return(&cloudformation.DescribeStackSetOutput{
			StackSet: &types.StackSet{
				StackSetId:   aws.String("arn:aws:cloudformation:us-east-1:123456789012:stackset/stack-set-name/guid"),
				StackSetName: aws.String("stack-set-name"),
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.StackSet
		wantErr   bool
		expectErr string
	}{
		{
			name: "query stack-set resource",
			expected: []types.StackSet{
				{
					StackSetId:   aws.String("arn:aws:cloudformation:us-east-1:123456789012:stackset/stack-set-name/guid"),
					StackSetName: aws.String("stack-set-name"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqCloudformationAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("stack-set")

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

			if actual.Service != "cloudformation" {
				t.Errorf("expected cloudformation, but got %v", actual.Service)
			}
			if actual.Resource != "stack-set" {
				t.Errorf("expected stack-set, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.StackSet)
				if !ok {
					t.Errorf("expected types.StackSet, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.StackSetId, tt.expected[i].StackSetId) {
					t.Errorf("expected %v, but got %v", tt.expected[i].StackSetId, actualOutput.StackSetId)
				}
				if !reflect.DeepEqual(actualOutput.StackSetName, tt.expected[i].StackSetName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].StackSetName, actualOutput.StackSetName)
				}
			}
		})
	}
}
