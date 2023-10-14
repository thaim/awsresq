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
			name:     "valid role resource",
			api:      AwsresqIamAPI{},
			resource: "role",
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
			name: "valid instance query",
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
