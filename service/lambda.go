package service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/rs/zerolog/log"
)

type awsLambdaAPI interface {
	ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

type AwsresqLambdaAPI struct {
	awsCfg aws.Config
	region []string
	apiClient map[string]awsLambdaAPI
}

func NewAwsresqLambdaAPI(c aws.Config, region []string) *AwsresqLambdaAPI {
	return &AwsresqLambdaAPI{
		awsCfg: c,
		region: region,
		apiClient: make(map[string]awsLambdaAPI, len(region)),
	}
}

func (api AwsresqLambdaAPI) Validate(resource string) bool {
	switch resource {
	case "function":
		return true
	}

	return false
}

func (api AwsresqLambdaAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "lambda",
		Resource: resource,
	}

	awsAPI := lambda.NewFromConfig(api.awsCfg)
	switch resource {
	case "function":
		listOutput, err := awsAPI.ListFunctions(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		for _, f := range listOutput.Functions {
			resultList.Results = append(resultList.Results, f)
		}
	default:
		log.Error().Msgf("resource '%s' not supported in lambda service", resource)
		return nil, fmt.Errorf("resource '%s' not supported in lambda service", resource)
	}

	return resultList, nil
}
