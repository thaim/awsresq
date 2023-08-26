package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/rs/zerolog/log"
)

type AwsresqLogsAPI struct {
	awsCfg aws.Config
	region []string
}

func NewAwsresqLogsAPI(c aws.Config, region []string) *AwsresqLogsAPI {
	return &AwsresqLogsAPI{
		awsCfg: c,
		region: region,
	}
}

func (api AwsresqLogsAPI) Validate(resource string) bool {
	switch resource {
	case "log-group":
		return true
	}

	return false
}

func (api AwsresqLogsAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service: "logs",
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
		Service: "logs",
		Resource: "log-group",
	}

	awsAPI := cloudwatchlogs.NewFromConfig(api.awsCfg, func(o *cloudwatchlogs.Options) {
		o.Region = r
	})
	listOutput, err := awsAPI.DescribeLogGroups(ctx, nil)
	if err != nil {
		log.Error().Msgf("error querying log groups in region %s: %s", r, err)
		return
	}

	for _, lg := range listOutput.LogGroups {
		resultList.Results = append(resultList.Results, lg)
	}

	ch <- resultList
}
