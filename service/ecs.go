//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsEcsAPI interface {
	ListTaskDefinitions(ctx context.Context, params *ecs.ListTaskDefinitionsInput, optFns ...func(*ecs.Options)) (*ecs.ListTaskDefinitionsOutput, error)
	DescribeTaskDefinition(ctx context.Context, params *ecs.DescribeTaskDefinitionInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTaskDefinitionOutput, error)
}

type AwsresqEcsAPI struct {
	awsCfg    aws.Config
	region    []string
	apiClient map[string]awsEcsAPI
}

func NewAwsresqEcsAPI(c aws.Config, region []string) AwsresqEcsAPI {
	return AwsresqEcsAPI{
		awsCfg:    c,
		region:    region,
		apiClient: make(map[string]awsEcsAPI, len(region)),
	}
}

func (api AwsresqEcsAPI) Validate(resource string) bool {
	validResources := []string{
		"service",
		"task-definition",
	}

	return slices.Contains(validResources, resource)
}

func (api AwsresqEcsAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "ecs",
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

		for range api.region {
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

func (api *AwsresqEcsAPI) queryTaskDefinition(ctx context.Context, ch chan ResultList, r string) {
	resultList := ResultList{
		Service:  "ecs",
		Resource: "task-definition",
	}

	if api.apiClient[r] == nil {
		api.apiClient[r] = ecs.NewFromConfig(api.awsCfg, func(o *ecs.Options) {
			o.Region = r
		})
	}

	listOutput, err := api.apiClient[r].ListTaskDefinitions(ctx, nil)
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
		output, err := api.apiClient[r].DescribeTaskDefinition(ctx, input)
		if err != nil {
			fmt.Println(err)
			return
		}

		resultList.Results = append(resultList.Results, output)
	}

	ch <- resultList
}
