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

	awsAPI := ecs.NewFromConfig(api.awsCfg)

	switch resource {
	case "task-definition":
		listOutput, err := awsAPI.ListTaskDefinitions(context.Background(), nil)
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
	default:
		log.Error().Msgf("resource '%s' not supported in ecs service", resource)
		return nil, fmt.Errorf("resource '%s' not supported in ecs service", resource)
	}

	return resultList, nil
}
