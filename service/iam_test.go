package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestIamValidate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqIamAPI
		resource string
		expect   bool
	}{
		{
			name:     "valid access-key resource",
			api:      AwsresqIamAPI{},
			resource: "access-key",
			expect:   true,
		},
		{
			name:     "valid group resource",
			api:      AwsresqIamAPI{},
			resource: "group",
			expect:   true,
		},
		{
			name:     "valid role resource",
			api:      AwsresqIamAPI{},
			resource: "role",
			expect:   true,
		},
		{
			name:     "valid user resource",
			api:      AwsresqIamAPI{},
			resource: "user",
			expect:   true,
		},
		{
			name:     "undefined resource",
			api:      AwsresqIamAPI{},
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

func TestIamAccessKeysQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsIamAPI(ctrl)

	mc.EXPECT().
		ListAccessKeys(gomock.Any(), nil).
		Return(&iam.ListAccessKeysOutput{
			AccessKeyMetadata: []types.AccessKeyMetadata{
				{
					AccessKeyId: aws.String("test-access-key"),
					UserName:    aws.String("test-user"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.AccessKeyMetadata
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid access-key query",
			expected: []types.AccessKeyMetadata{
				{
					AccessKeyId:  aws.String("test-access-key"),
					UserName:     aws.String("test-user"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqIamAPI(config, []string{"us-east-1"})
			api.apiClient["us-east-1"] = mc

			actual, err := api.Query("access-key")

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

			if actual.Service != "iam" {
				t.Errorf("expected iam, but got %v", actual.Service)
			}
			if actual.Resource != "access-key" {
				t.Errorf("expected access-key, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.AccessKeyMetadata)
				if !ok {
					t.Errorf("expected types.AccessKeyMetadata, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.AccessKeyId, tt.expected[i].AccessKeyId) {
					t.Errorf("expected %v, but got %v", tt.expected[i].AccessKeyId, actualOutput.AccessKeyId)
				}
				if !reflect.DeepEqual(actualOutput.UserName, tt.expected[i].UserName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].UserName, actualOutput.UserName)
				}
			}
		})
	}
}

func TestIamGroupQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsIamAPI(ctrl)

	mc.EXPECT().
		ListGroups(gomock.Any(), nil).
		Return(&iam.ListGroupsOutput{
			Groups: []types.Group{
				{
					GroupName: aws.String("test-group"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.Group
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid group query",
			expected: []types.Group{
				{
					GroupName: aws.String("test-group"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqIamAPI(config, []string{"us-east-1"})
			api.apiClient["us-east-1"] = mc

			actual, err := api.Query("group")

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

			if actual.Service != "iam" {
				t.Errorf("expected iam, but got %v", actual.Service)
			}
			if actual.Resource != "group" {
				t.Errorf("expected group, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Group)
				if !ok {
					t.Errorf("expected types.Group, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.GroupName, tt.expected[i].GroupName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].GroupName, actualOutput.GroupName)
				}
			}
		})
	}
}

func TestIamPolicyQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsIamAPI(ctrl)

	mc.EXPECT().
		ListPolicies(gomock.Any(), &iam.ListPoliciesInput{
			Scope: types.PolicyScopeTypeLocal,
		}).
		Return(&iam.ListPoliciesOutput{
			Policies: []types.Policy{
				{
					PolicyName: aws.String("test-policy"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.Policy
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid policy query",
			expected: []types.Policy{
				{
					PolicyName: aws.String("test-policy"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqIamAPI(config, []string{"us-east-1"})
			api.apiClient["us-east-1"] = mc

			actual, err := api.Query("policy")

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

			if actual.Service != "iam" {
				t.Errorf("expected iam, but got %v", actual.Service)
			}
			if actual.Resource != "policy" {
				t.Errorf("expected policy, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Policy)
				if !ok {
					t.Errorf("expected types.Group, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.PolicyName, tt.expected[i].PolicyName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].PolicyName, actualOutput.PolicyName)
				}
			}
		})
	}
}

func TestIamRoleQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsIamAPI(ctrl)

	mc.EXPECT().
		ListRoles(gomock.Any(), nil).
		Return(&iam.ListRolesOutput{
			Roles: []types.Role{
				{
					RoleName: aws.String("test-role"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.Role
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid role query",
			expected: []types.Role{
				{
					RoleName: aws.String("test-role"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqIamAPI(config, []string{"us-east-1"})
			api.apiClient["us-east-1"] = mc

			actual, err := api.Query("role")

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

			if actual.Service != "iam" {
				t.Errorf("expected iam, but got %v", actual.Service)
			}
			if actual.Resource != "role" {
				t.Errorf("expected role, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.Role)
				if !ok {
					t.Errorf("expected types.Role, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.RoleName, tt.expected[i].RoleName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].RoleName, actualOutput.RoleName)
				}
			}
		})
	}
}

func TestIamUserQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsIamAPI(ctrl)

	mc.EXPECT().
		ListUsers(gomock.Any(), nil).
		Return(&iam.ListUsersOutput{
			Users: []types.User{
				{
					UserName: aws.String("test-user"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.User
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid user query",
			expected: []types.User{
				{
					UserName: aws.String("test-user"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqIamAPI(config, []string{"us-east-1"})
			api.apiClient["us-east-1"] = mc

			actual, err := api.Query("user")

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

			if actual.Service != "iam" {
				t.Errorf("expected iam, but got %v", actual.Service)
			}
			if actual.Resource != "user" {
				t.Errorf("expected user, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.User)
				if !ok {
					t.Errorf("expected types.User, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.UserName, tt.expected[i].UserName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].UserName, actualOutput.UserName)
				}
			}
		})
	}
}
