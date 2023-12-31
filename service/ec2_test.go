package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestEc2Validate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqEc2API
		resource string
		expect   bool
	}{
		{
			name:     "valid instance resource",
			api:      AwsresqEc2API{},
			resource: "instance",
			expect:   true,
		},
		{
			name:     "undefined resource",
			api:      AwsresqEc2API{},
			resource: "undefined",
			expect:   false,
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

func TestEc2InstanceQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsEc2API(ctrl)

	mc.EXPECT().
		DescribeInstances(gomock.Any(), nil).
		Return(&ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							InstanceId:   aws.String("i-1234567890abcdef0"),
							InstanceType: types.InstanceTypeT2Micro,
						},
					},
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.Instance
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid instance query",
			expected: []types.Instance{
				{
					InstanceId:   aws.String("i-1234567890abcdef0"),
					InstanceType: types.InstanceTypeT2Micro,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqEc2API(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("instance")

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

			if actual.Service != "ec2" {
				t.Errorf("expected ec2, but got %v", actual.Service)
			}
			if actual.Resource != "instance" {
				t.Errorf("expected instance, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Instance)
				if !ok {
					t.Errorf("expected types.Cluster, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.InstanceId, tt.expected[i].InstanceId) {
					t.Errorf("expected %v, but got %v", tt.expected[i].InstanceId, actualOutput.InstanceId)
				}
				if !reflect.DeepEqual(actualOutput.InstanceType, tt.expected[i].InstanceType) {
					t.Errorf("expected %v, but got %v", tt.expected[i].InstanceType, actualOutput.InstanceType)
				}
			}
		})
	}
}

func TestEc2SecurityGroupQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsEc2API(ctrl)

	mc.EXPECT().
		DescribeSecurityGroups(gomock.Any(), nil).
		Return(&ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: []types.SecurityGroup{
				{
					GroupId:   aws.String("sg-1234567890abcdef0"),
					GroupName: aws.String("test"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.SecurityGroup
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid security-group query",
			expected: []types.SecurityGroup{
				{
					GroupId:   aws.String("sg-1234567890abcdef0"),
					GroupName: aws.String("test"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqEc2API(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("security-group")

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

			if actual.Service != "ec2" {
				t.Errorf("expected ec2, but got %v", actual.Service)
			}
			if actual.Resource != "security-group" {
				t.Errorf("expected security-group, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.SecurityGroup)
				if !ok {
					t.Errorf("expected types.SecurityGroup, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.GroupId, tt.expected[i].GroupId) {
					t.Errorf("expected %v, but got %v", tt.expected[i].GroupId, actualOutput.GroupId)
				}
				if !reflect.DeepEqual(actualOutput.GroupName, tt.expected[i].GroupName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].GroupName, actualOutput.GroupName)
				}
			}
		})
	}
}

func TestEc2VpcQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsEc2API(ctrl)

	mc.EXPECT().
		DescribeVpcs(gomock.Any(), nil).
		Return(&ec2.DescribeVpcsOutput{
			Vpcs: []types.Vpc{
				{
					VpcId:     aws.String("vpc-1234567890abcdef0"),
					CidrBlock: aws.String("172.31.0.0/16"),
					IsDefault: aws.Bool(true),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.Vpc
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid vpc query",
			expected: []types.Vpc{
				{
					VpcId:     aws.String("vpc-1234567890abcdef0"),
					CidrBlock: aws.String("172.31.0.0/16"),
					IsDefault: aws.Bool(true),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqEc2API(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("vpc")

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

			if actual.Service != "ec2" {
				t.Errorf("expected ec2, but got %v", actual.Service)
			}
			if actual.Resource != "vpc" {
				t.Errorf("expected vpc, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Vpc)
				if !ok {
					t.Errorf("expected types.Vpc, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.VpcId, tt.expected[i].VpcId) {
					t.Errorf("expected %v, but got %v", tt.expected[i].VpcId, actualOutput.VpcId)
				}
				if !reflect.DeepEqual(actualOutput.CidrBlock, tt.expected[i].CidrBlock) {
					t.Errorf("expected %v, but got %v", tt.expected[i].CidrBlock, actualOutput.CidrBlock)
				}
				if !reflect.DeepEqual(actualOutput.IsDefault, tt.expected[i].IsDefault) {
					t.Errorf("expected %v, but got %v", tt.expected[i].IsDefault, actualOutput.IsDefault)
				}
			}
		})
	}
}
