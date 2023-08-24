package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/thaim/awsresq/mock"
	"github.com/golang/mock/gomock"
)


func TestLambdaValidate(t *testing.T) {
	cases := []struct {
		name string
		api AwsresqLambdaAPI
		resource string
		expected bool
	}{
		{
			name: "validate function resource",
			api: AwsresqLambdaAPI{},
			resource: "function",
			expected: true,
		},
		{
			name: "validate undefined resource",
			api: AwsresqLambdaAPI{},
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

func TestLambdaQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mock_service.NewMockawsLambdaAPI(ctrl)

	mc.EXPECT().
		ListFunctions(gomock.Any(), nil).
		Return(&lambda.ListFunctionsOutput{
			Functions: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("testapp"),
				},
			},
		}, nil).
		AnyTimes()

	cases := []struct {
		name string
		expected *lambda.ListFunctionsOutput
		wantErr bool
		expectErr string
	}{
		{
			name: "query function resource",
			expected: &lambda.ListFunctionsOutput{
				Functions: []types.FunctionConfiguration{
					{
						FunctionName: aws.String("testapp"),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config, _ := config.LoadDefaultConfig(context.TODO())
			api := NewAwsresqLambdaAPI(config, []string{"ap-northeast-1"})
			api.apiClient["ap-northeast-1"] = mc

			actual, err := api.Query("function")

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
				if err.Error() != tt.expectErr {
					t.Errorf("expected error message %v, but got %v", tt.expectErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}

			if actual.Service != "lambda" {
				t.Errorf("expected service %v, but got %v", "lambda", actual.Service)
			}
			if actual.Resource != "function" {
				t.Errorf("expected resource 'function', but got %v", actual.Resource)
			}

			if len(tt.expected.Functions) != len(actual.Results) {
				t.Errorf("expected %v results, but got %v", len(tt.expected.Functions), len(actual.Results))
			}
			for i := range tt.expected.Functions {
				actualOutput, ok := actual.Results[i].(types.FunctionConfiguration)
				if !ok {
					t.Errorf("expected type type.FunctionConfiguration, but got %T", actual.Results[i])
				}
				if !reflect.DeepEqual(tt.expected.Functions[i], actualOutput) {
					t.Errorf("expected %+v, but got %+v", tt.expected.Functions[i], actualOutput)
				}
			}
		})
	}
}
