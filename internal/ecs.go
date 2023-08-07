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
	region []string
}

func NewAwsEcsAPI(c aws.Config, region []string) *AwsEcsAPI {
	return &AwsEcsAPI{
		awsCfg: c,
		region: region,
	}
}

func (api *AwsEcsAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service: "ecs",
		Resource: resource,
	}
	var err error = nil

	switch resource {
	case "task-definition":
		ch := make(chan ResultList)

		for _, r := range api.region {
			go api.queryTaskDefinition(ch, r)
		}

		for _ = range api.region {
			result := <-ch
			if result.Results != nil {
				resultList.Results = append(resultList.Results, result.Results...)
			}
		}

	default:
		log.Error().Msgf("resource '%s' not supported in ecs service", resource)
		return nil, fmt.Errorf("resource '%s' not supported in ecs service", resource)
	}

	return resultList, err
}

func (api *AwsEcsAPI) queryTaskDefinition(ch chan ResultList, r string) {
	resultList := ResultList{
		Service: "ecs",
		Resource: "task-definition",
	}

	awsApi := ecs.NewFromConfig(api.awsCfg, func(o *ecs.Options) {
		o.Region = r
	})

	listOutput, err := awsApi.ListTaskDefinitions(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, arn := range listOutput.TaskDefinitionArns {
		input := &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String(arn),
			Include: []types.TaskDefinitionField{
				types.TaskDefinitionFieldTags,
			},
		}
		output, err := awsApi.DescribeTaskDefinition(context.Background(), input)
		if err != nil {
			fmt.Println(err)
			return
		}

		resultList.Results = append(resultList.Results, output)
	}

	ch <- resultList
}
