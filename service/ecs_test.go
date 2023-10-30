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
			name:     "validate cluster resource",
			api:      AwsresqEcsAPI{},
			resource: "cluster",
			expected: true,
		},
		{
			name:     "validate service resource",
			api:      AwsresqEcsAPI{},
			resource: "service",
			expected: true,
		},
		{
			name:     "validate task resource",
			api:      AwsresqEcsAPI{},
			resource: "task",
			expected: true,
		},
		{
			name:     "validate task-definition resource",
			api:      AwsresqEcsAPI{},
			resource: "task-definition",
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

func TestEcsClusterQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsEcsAPI(ctrl)

	mc.EXPECT().
		ListClusters(gomock.Any(), nil).
		Return(&ecs.ListClustersOutput{
			ClusterArns: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:cluster/test-cluster",
			},
		}, nil).
		AnyTimes()
	mc.EXPECT().
		DescribeClusters(gomock.Any(), &ecs.DescribeClustersInput{
			Clusters: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:cluster/test-cluster",
			},
			Include: []types.ClusterField{
				types.ClusterFieldTags,
				types.ClusterFieldStatistics,
				types.ClusterFieldSettings,
				types.ClusterFieldConfigurations,
				types.ClusterFieldAttachments,
			},
		}).
		Return(&ecs.DescribeClustersOutput{
			Clusters: []types.Cluster{
				{
					ClusterArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/test-cluster"),
					ClusterName: aws.String("test-cluster"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.Cluster
		wantErr   bool
		expectErr string
	}{
		{
			name: "query cluster resource",
			expected: []types.Cluster{
				{
					ClusterArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/test-cluster"),
					ClusterName: aws.String("test-cluster"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqEcsAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("cluster")

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

			if actual.Service != "ecs" {
				t.Errorf("expected ecs, but got %v", actual.Service)
			}
			if actual.Resource != "cluster" {
				t.Errorf("expected cluster, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Cluster)
				if !ok {
					t.Errorf("expected types.Cluster, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.ClusterArn, tt.expected[i].ClusterArn) {
					t.Errorf("expected %v, but got %v", tt.expected[i].ClusterArn, actualOutput.ClusterArn)
				}
				if !reflect.DeepEqual(actualOutput.ClusterName, tt.expected[i].ClusterName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].ClusterName, actualOutput.ClusterName)
				}
			}
		})
	}
}

func TestEcsTaskQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsEcsAPI(ctrl)

	mc.EXPECT().
		ListClusters(gomock.Any(), nil).
		Return(&ecs.ListClustersOutput{
			ClusterArns: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster01",
				"arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster02",
			},
		}, nil).
		AnyTimes()

	mc.EXPECT().
		ListTasks(gomock.Any(), &ecs.ListTasksInput{
			Cluster: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster01"),
		}).
		Return(&ecs.ListTasksOutput{
			TaskArns: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:task/testcluster01/74de0355a10a4f979ac495c14EXAMPLE",
				"arn:aws:ecs:ap-northeast-1:012345678901:task/testcluster01/d789e94343414c25b9f6bd59eEXAMPLE",
			},
		}, nil).
		AnyTimes()
	mc.EXPECT().
		ListTasks(gomock.Any(), &ecs.ListTasksInput{
			Cluster: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster02"),
		}).
		Return(&ecs.ListTasksOutput{
			TaskArns: []string{},
		}, nil).
		AnyTimes()

	mc.EXPECT().
		DescribeTasks(gomock.Any(), &ecs.DescribeTasksInput{
			Cluster: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster01"),
			Tasks: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:task/testcluster01/74de0355a10a4f979ac495c14EXAMPLE",
			},
		}).
		Return(&ecs.DescribeTasksOutput{
			Tasks: []types.Task{
				{
					TaskArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task/testcluster01/74de0355a10a4f979ac495c14EXAMPLE"),
				},
			},
		}, nil).
		AnyTimes()
	mc.EXPECT().
		DescribeTasks(gomock.Any(), &ecs.DescribeTasksInput{
			Cluster: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster01"),
			Tasks: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:task/testcluster01/d789e94343414c25b9f6bd59eEXAMPLE",
			},
		}).
		Return(&ecs.DescribeTasksOutput{
			Tasks: []types.Task{
				{
					TaskArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task/testcluster01/d789e94343414c25b9f6bd59eEXAMPLE"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		resource  string
		expected  []*types.Task
		wantErr   bool
		expectErr string
	}{
		{
			name:     "query task resource",
			resource: "task",
			expected: []*types.Task{
				{
					TaskArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task/testcluster01/74de0355a10a4f979ac495c14EXAMPLE"),
				},
				{
					TaskArn: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:task/testcluster01/d789e94343414c25b9f6bd59eEXAMPLE"),
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
				t.Errorf("expected resource 'task', but got '%v'", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %d results, but got %d", len(tt.expected), len(actual.Results))
			}
			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(*types.Task)
				if !ok {
					t.Errorf("expected type *ecs.Task, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(tt.expected[i].TaskArn, actualOutput.TaskArn) {
					t.Errorf("expected %+v, but got %+v", tt.expected[i].TaskArn, actualOutput.TaskArn)
				}
			}
		})
	}
}

func TestEcsTaskDefinitionQuery(t *testing.T) {
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

func TestEcsServiceQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsEcsAPI(ctrl)

	mc.EXPECT().
		ListClusters(gomock.Any(), nil).
		Return(&ecs.ListClustersOutput{
			ClusterArns: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster01",
				"arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster02",
			},
		}, nil).
		AnyTimes()

	mc.EXPECT().
		ListServices(gomock.Any(), &ecs.ListServicesInput{
			Cluster: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster01"),
		}).
		Return(&ecs.ListServicesOutput{
			ServiceArns: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:service/testcluster01/testservice01",
			},
		}, nil).
		AnyTimes()
	mc.EXPECT().
		ListServices(gomock.Any(), &ecs.ListServicesInput{
			Cluster: aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster02"),
		}).
		Return(&ecs.ListServicesOutput{
			ServiceArns: []string{
				"arn:aws:ecs:ap-northeast-1:012345678901:service/testcluster02/testservice02",
			},
		}, nil).
		AnyTimes()

	mc.EXPECT().
		DescribeServices(gomock.Any(), &ecs.DescribeServicesInput{
			Cluster:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster01"),
			Services: []string{"arn:aws:ecs:ap-northeast-1:012345678901:service/testcluster01/testservice01"},
		}).
		Return(&ecs.DescribeServicesOutput{
			Services: []types.Service{
				{
					ClusterArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster01"),
					ServiceArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:service/testcluster01/testservice01"),
					ServiceName: aws.String("testservice01"),
				},
			},
		}, nil).
		AnyTimes()
	mc.EXPECT().
		DescribeServices(gomock.Any(), &ecs.DescribeServicesInput{
			Cluster:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster02"),
			Services: []string{"arn:aws:ecs:ap-northeast-1:012345678901:service/testcluster02/testservice02"},
		}).
		Return(&ecs.DescribeServicesOutput{
			Services: []types.Service{
				{
					ClusterArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster02"),
					ServiceArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:service/testcluster02/testservice02"),
					ServiceName: aws.String("testservice02"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []*types.Service
		wantErr   bool
		expectErr string
	}{
		{
			name: "query service resource",
			expected: []*types.Service{
				{
					ClusterArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster01"),
					ServiceArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:service/testcluster01/testservice01"),
					ServiceName: aws.String("testservice01"),
				},
				{
					ClusterArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:cluster/testcluster02"),
					ServiceArn:  aws.String("arn:aws:ecs:ap-northeast-1:012345678901:service/testcluster02/testservice02"),
					ServiceName: aws.String("testservice02"),
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqEcsAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("service")

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
			if actual.Resource != "service" {
				t.Errorf("expected resource 'task-definition', but got '%v'", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %d results, but got %d", len(tt.expected), len(actual.Results))
			}
			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Service)
				if !ok {
					t.Errorf("expected type types.Service, but got %T", actual.Results[i])
				}
				if *tt.expected[i].ClusterArn != *actualOutput.ClusterArn {
					t.Errorf("expected %+v, but got %+v", tt.expected[i].ClusterArn, actualOutput.ClusterArn)
				}
				if *tt.expected[i].ServiceArn != *actualOutput.ServiceArn {
					t.Errorf("expected %+v, but got %+v", tt.expected[i].ServiceArn, actualOutput.ServiceArn)
				}
				if *tt.expected[i].ServiceName != *actualOutput.ServiceName {
					t.Errorf("expected %+v, but got %+v", *tt.expected[i].ServiceName, *actualOutput.ServiceName)
				}
			}
		})
	}
}
