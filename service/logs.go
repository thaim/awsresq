//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsLogsAPI interface {
	DescribeLogGroups(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
}

type AwsresqLogsAPI struct {
	awsCfg    aws.Config
	region    []string
	apiClient map[string]awsLogsAPI
}

func NewAwsresqLogsAPI(c aws.Config, region []string) *AwsresqLogsAPI {
	return &AwsresqLogsAPI{
		awsCfg:    c,
		region:    region,
		apiClient: make(map[string]awsLogsAPI, len(region)),
	}
}

func (api AwsresqLogsAPI) Validate(resource string) bool {
	validResoruces := []string{
		"log-group",
	}

	return slices.Contains(validResoruces, resource)
}

func (api AwsresqLogsAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "logs",
		Resource: resource,
	}

	switch resource {
	case "log-group":
		ch := make(chan ResultList)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for _, r := range api.region {
			go api.queryLogGroup(ctx, ch, r)
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
	default:
		log.Error().Msgf("resource '%s' not supported in logs service", resource)
		return nil, fmt.Errorf("resource '%s' not supported in logs service", resource)
	}

	return resultList, nil
}

func (api *AwsresqLogsAPI) queryLogGroup(ctx context.Context, ch chan ResultList, r string) {
	resultList := ResultList{
		Service:  "logs",
		Resource: "log-group",
	}

	if api.apiClient[r] == nil {
		api.apiClient[r] = cloudwatchlogs.NewFromConfig(api.awsCfg, func(o *cloudwatchlogs.Options) {
			o.Region = r
		})
	}

	listOutput, err := api.apiClient[r].DescribeLogGroups(ctx, nil)
	if err != nil {
		log.Error().Msgf("error querying log groups in region %s: %s", r, err)
		return
	}

	for _, lg := range listOutput.LogGroups {
		resultList.Results = append(resultList.Results, lg)
	}

	ch <- resultList
}
