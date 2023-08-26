package service

import (
	"context"
	"fmt"

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
		for _, r := range api.region {
			awsAPI := cloudwatchlogs.NewFromConfig(api.awsCfg, func(o *cloudwatchlogs.Options) {
				o.Region = r
			})
			listOutput, err := awsAPI.DescribeLogGroups(context.Background(), nil)
			if err != nil {
				return nil, err
			}
			for _, lg := range listOutput.LogGroups {
				resultList.Results = append(resultList.Results, lg)
			}
		}
	default:
		log.Error().Msgf("resource '%s' not supported in logs service", resource)
		return nil, fmt.Errorf("resource '%s' not supported in logs service", resource)
	}

	return resultList, nil
}
