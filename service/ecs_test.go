package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestEcsValidate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqEcsAPI
		resource string
		expected bool
	}{
		{
			name:     "validate task-definition resource",
			api:      AwsresqEcsAPI{},
			resource: "task-definition",
			expected: true,
		},
		{
			name:     "validate service resource",
			api:      AwsresqEcsAPI{},
			resource: "service",
			expected: true,
		},
		{
			name:     "validate undefined resource",
			api:      AwsresqEcsAPI{},
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

func TestEcsQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsEcsAPI(ctrl)

	mc.EXPECT().
		ListTaskDefinitions(gomock.Any(), nil).
		Return(&ecs.ListTaskDefinitionsOutput{
			TaskDefinitionArns: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:task-definition/testapp:1",
				"arn:aws:ecs:ap-northeast-1:012345678901:task-definition/testapp:2",
				"arn:aws:ecs:ap-northeast-1:012345678901:task-definition/sampleapp:1",
			},
		}, nil).
		AnyTimes()

	mc.EXPECT().
		DescribeTaskDefinition(gomock.Any(), &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task-definition/testapp:1"),
			Include: []types.TaskDefinitionField{
				types.TaskDefinitionFieldTags,
			},
		}).
		Return(&ecs.DescribeTaskDefinitionOutput{
			Tags: []types.Tag{},
			TaskDefinition: &types.TaskDefinition{
				TaskDefinitionArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task-definition/testapp:1"),
			},
		}, nil).
		AnyTimes()
	mc.EXPECT().
		DescribeTaskDefinition(gomock.Any(), &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task-definition/testapp:2"),
			Include: []types.TaskDefinitionField{
				types.TaskDefinitionFieldTags,
			},
		}).
		Return(&ecs.DescribeTaskDefinitionOutput{
			Tags: []types.Tag{},
			TaskDefinition: &types.TaskDefinition{
				TaskDefinitionArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task-definition/testapp:2"),
			},
		}, nil).
		AnyTimes()
	mc.EXPECT().
		DescribeTaskDefinition(gomock.Any(), &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task-definition/sampleapp:1"),
			Include: []types.TaskDefinitionField{
				types.TaskDefinitionFieldTags,
			},
		}).
		Return(&ecs.DescribeTaskDefinitionOutput{
			Tags: []types.Tag{},
			TaskDefinition: &types.TaskDefinition{
				TaskDefinitionArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task-definition/sampleapp:1"),
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		resource  string
		expected  []*types.TaskDefinition
		wantErr   bool
		expectErr string
	}{
		{
			name:     "query task-definition resource",
			resource: "task-definition",
			expected: []*types.TaskDefinition{
				{
					TaskDefinitionArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task-definition/testapp:1"),
				},
				{
					TaskDefinitionArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task-definition/testapp:2"),
				},
				{
					TaskDefinitionArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task-definition/sampleapp:1"),
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqEcsAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query(tt.resource)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error '%s', but got no error", tt.expectErr)
				} else if !strings.Contains(err.Error(), tt.expectErr) {
					t.Errorf("expected error '%s', but got '%s'", tt.expectErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if actual.Service != "ecs" {
				t.Errorf("expected service 'ecs', but got '%v'", actual.Service)
			}
			if actual.Resource != "task-definition" {
				t.Errorf("expected resource 'task-definition', but got '%v'", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %d results, but got %d", len(tt.expected), len(actual.Results))
			}
			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(*types.TaskDefinition)
				if !ok {
					t.Errorf("expected type *ecs.DescribeTaskDefinitionOutput, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(tt.expected[i], actualOutput) {
					t.Errorf("expected %+v, but got %+v", tt.expected[i], actualOutput)
				}
			}
		})
	}
}
