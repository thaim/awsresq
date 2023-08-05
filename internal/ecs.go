package internal

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/rs/zerolog/log"
)

type AwsEcsAPI struct {
	awsCfg aws.Config
	awsApi *ecs.Client
}

func NewAwsEcsAPI(c aws.Config) *AwsEcsAPI {
	return &AwsEcsAPI{
		awsCfg: c,
	}
}

func (api *AwsEcsAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service: "ecs",
		Resource: resource,
	}
	var err error = nil

	api.awsApi = ecs.NewFromConfig(api.awsCfg)

	switch resource {
	case "task-definition":
		resultList, err = api.queryTaskDefinition()
	default:
		log.Error().Msgf("resource '%s' not supported in ecs service", resource)
		return nil, fmt.Errorf("resource '%s' not supported in ecs service", resource)
	}

	return resultList, err
}

func (api *AwsEcsAPI) queryTaskDefinition() (*ResultList, error) {
	resultList := &ResultList{
		Service: "ecs",
		Resource: "task-definition",
	}

	awsAPI := ecs.NewFromConfig(api.awsCfg)

	listOutput, err := api.awsApi.ListTaskDefinitions(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	for _, arn := range listOutput.TaskDefinitionArns {
		input := &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String(arn),
			Include: []types.TaskDefinitionField{
				types.TaskDefinitionFieldTags,
			},
		}
		output, err := awsAPI.DescribeTaskDefinition(context.Background(), input)
		if err != nil {
			return nil, err
		}

		resultList.Results = append(resultList.Results, output)
	}

	return resultList, nil
}
