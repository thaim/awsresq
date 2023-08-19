package service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/rs/zerolog/log"
)

type AwsLambdaAPI struct {
	awsCfg aws.Config
}

func NewAwsLambdaAPI(c aws.Config) *AwsLambdaAPI {
	return &AwsLambdaAPI{
		awsCfg: c,
	}
}

func (api AwsLambdaAPI) Validate(resource string) bool {
	switch resource {
	case "function":
		return true
	}

	return false
}

func (api AwsLambdaAPI) Query(resource string) (*ResultList, error) {
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