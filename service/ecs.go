package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/rs/zerolog/log"
)

type AwsEcsAPI struct {
	awsCfg aws.Config
	region []string
}

func NewAwsEcsAPI(c aws.Config, region []string) AwsEcsAPI {
	return AwsEcsAPI{
		awsCfg: c,
		region: region,
	}
}

func (api AwsEcsAPI) Validate(resource string) bool {
	switch resource {
		case "task-definition":
		return true
	}

	return false
}

func (api AwsEcsAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service: "ecs",
		Resource: resource,
	}
	var err error = nil

	switch resource {
	case "task-definition":
		ch := make(chan ResultList)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		for _, r := range api.region {
			go api.queryTaskDefinition(ctx, ch, r)
		}

		for _ = range api.region {
			select {
			case result := <-ch:
				if result.Results != nil {
					resultList.Results = append(resultList.Results, result.Results...)
				}
			case <-ctx.Done():
				return resultList, ctx.Err()
			}
		}

	default:
		log.Error().Msgf("resource '%s' not supported in ecs service", resource)
		return nil, fmt.Errorf("resource '%s' not supported in ecs service", resource)
	}

	return resultList, err
}

func (api *AwsEcsAPI) queryTaskDefinition(ctx context.Context, ch chan ResultList, r string) {
	resultList := ResultList{
		Service: "ecs",
		Resource: "task-definition",
	}

	awsApi := ecs.NewFromConfig(api.awsCfg, func(o *ecs.Options) {
		o.Region = r
	})

	listOutput, err := awsApi.ListTaskDefinitions(ctx, nil)
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
		output, err := awsApi.DescribeTaskDefinition(ctx, input)
		if err != nil {
			fmt.Println(err)
			return
		}

		resultList.Results = append(resultList.Results, output)
	}

	ch <- resultList
}
