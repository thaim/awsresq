package service

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/golang/mock/gomock"
	"github.com/thaim/awsresq/mock"
)

func TestRoute53Validate(t *testing.T) {
	cases := []struct {
		name     string
		api      AwsresqRoute53API
		resource string
		expected bool
	}{
		{
			name:     "validate hosted-zone resource",
			api:      AwsresqRoute53API{},
			resource: "hosted-zone",
			expected: true,
		},
		{
			name:     "validate undefined resource",
			api:      AwsresqRoute53API{},
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

func TestRoute53HostedZoneQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsRoute53API(ctrl)

	mc.EXPECT().
		ListHostedZones(gomock.Any(), nil).
		Return(&route53.ListHostedZonesOutput{
			HostedZones: []types.HostedZone{
				{
					Id:   aws.String("/hostedzone/123456789012"),
					Name: aws.String("example.com."),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name      string
		expected  []types.HostedZone
		wantErr   bool
		expectErr string
	}{
		{
			name: "query hosted-zone resource",
			expected: []types.HostedZone{
				{
					Id:   aws.String("/hostedzone/123456789012"),
					Name: aws.String("example.com."),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqRoute53API(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("hosted-zone")

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

			if actual.Service != "route53" {
				t.Errorf("expected route53, but got %v", actual.Service)
			}
			if actual.Resource != "hosted-zone" {
				t.Errorf("expected hosted-zone, but got %v", actual.Resource)
			}

			if len(tt.expected) != len(actual.Results) {
				t.Errorf("expected %v, but got %v", len(tt.expected), len(actual.Results))
			}

			for i := range tt.expected {
				actualOutput, ok := actual.Results[i].(types.HostedZone)
				if !ok {
					t.Errorf("expected types.Stack, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(actualOutput.Id, tt.expected[i].Id) {
					t.Errorf("expected %v, but got %v", tt.expected[i].Id, actualOutput.Id)
				}
				if !reflect.DeepEqual(actualOutput.Name, tt.expected[i].Name) {
					t.Errorf("expected %v, but got %v", tt.expected[i].Name, actualOutput.Name)
				}
			}
		})
	}
}
