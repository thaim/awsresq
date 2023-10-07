//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsCloudformationAPI interface {
	ListStackSets(ctx context.Context, params *cloudformation.ListStackSetsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListStackSetsOutput, error)
	DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
	DescribeStackSet(ctx context.Context, params *cloudformation.DescribeStackSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackSetOutput, error)
}

type AwsresqCloudformationAPI struct {
	awsCfg    aws.Config
	region    []string
	apiClient map[string]awsCloudformationAPI
}

func NewAwsresqCloudformationAPI(c aws.Config, region []string) *AwsresqCloudformationAPI {
	return &AwsresqCloudformationAPI{
		awsCfg:    c,
		region:    region,
		apiClient: make(map[string]awsCloudformationAPI, len(region)),
	}
}

func (api AwsresqCloudformationAPI) Validate(resource string) bool {
	validResoruces := []string{
		"stack",
		"stack-set",
	}

	return slices.Contains(validResoruces, resource)
}

func (api AwsresqCloudformationAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "cloudformation",
		Resource: resource,
	}

	var apiQuery ResourceQueryAPI
	switch resource {
	case "stack":
		apiQuery = api.queryCloudformationStack
	case "stack-set":
		apiQuery = api.queryCloudformationStackSet
	default:
		return nil, fmt.Errorf("resource %s not supported in cloudformation service", resource)
	}

	ch := make(chan ResultList)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, r := range api.region {
		go apiQuery(ctx, ch, r)
	}

	for range api.region {
		select {
		case result := <-ch:
			if result.Results != nil {
				resultList.Results = append(resultList.Results, result.Results...)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return resultList, nil
}

func (api *AwsresqCloudformationAPI) queryCloudformationStack(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service:  "cloudformation",
		Resource: "stack",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = cloudformation.NewFromConfig(api.awsCfg, func(o *cloudformation.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].DescribeStacks(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("error querying cloudformation stacks for region %s", region)
		return
	}

	for _, stack := range listOutput.Stacks {
		resultList.Results = append(resultList.Results, stack)
	}

	ch <- resultList
}

func (api *AwsresqCloudformationAPI) queryCloudformationStackSet(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service:  "cloudformation",
		Resource: "stack-set",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = cloudformation.NewFromConfig(api.awsCfg, func(o *cloudformation.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].ListStackSets(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("error querying cloudformation stack-sets for region %s", region)
		return
	}

	for _, stackSet := range listOutput.Summaries {
		describeOutput, err := api.apiClient[region].DescribeStackSet(ctx, &cloudformation.DescribeStackSetInput{
			StackSetName: stackSet.StackSetName,
		})
		if err != nil {
			log.Error().Err(err).Msgf("error describing cloudformation stack-set %s for region %s", *stackSet.StackSetName, region)
			continue
		}

		resultList.Results = append(resultList.Results, *describeOutput.StackSet)
	}

	ch <- resultList
}
