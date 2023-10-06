package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestConfigValidate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqConfigAPI
		resource string
		expect   bool
	}{
		{
			name:     "valid rule resource",
			api:      AwsresqConfigAPI{},
			resource: "rule",
			expect:   true,
		},
		{
			name:     "undefined resource",
			api:      AwsresqConfigAPI{},
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

func TestConfigRuleQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsConfigAPI(ctrl)

	mc.EXPECT().
		DescribeConfigRules(gomock.Any(), nil).
		Return(&configservice.DescribeConfigRulesOutput{
			ConfigRules: []types.ConfigRule{
				{
					ConfigRuleArn:  aws.String("arn:aws:config:ap-northeast-1:123456789012:config-rule/config-rule-123456"),
					ConfigRuleId:   aws.String("config-rule-123456"),
					ConfigRuleName: aws.String("config-rule-123456"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.ConfigRule
		wantErr   bool
		expectErr string
	}{
		{
			name: "valid instance query",
			expected: []types.ConfigRule{
				{
					ConfigRuleArn:  aws.String("arn:aws:config:ap-northeast-1:123456789012:config-rule/config-rule-123456"),
					ConfigRuleId:   aws.String("config-rule-123456"),
					ConfigRuleName: aws.String("config-rule-123456"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqConfigAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("rule")

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

			if actual.Service != "config" {
				t.Errorf("expected config, but got %v", actual.Service)
			}
			if actual.Resource != "rule" {
				t.Errorf("expected rule, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.ConfigRule)
				if !ok {
					t.Errorf("expected types.ConfigRule, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.ConfigRuleArn, tt.expected[i].ConfigRuleArn) {
					t.Errorf("expected %v, but got %v", tt.expected[i].ConfigRuleArn, actualOutput.ConfigRuleArn)
				}
				if !reflect.DeepEqual(actualOutput.ConfigRuleId, tt.expected[i].ConfigRuleId) {
					t.Errorf("expected %v, but got %v", tt.expected[i].ConfigRuleId, actualOutput.ConfigRuleId)
				}
				if !reflect.DeepEqual(actualOutput.ConfigRuleName, tt.expected[i].ConfigRuleName) {
					t.Errorf("expected %v, but got %v", tt.expected[i].ConfigRuleName, actualOutput.ConfigRuleName)
				}
			}
		})
	}
}
