package service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/rs/zerolog/log"
)

type AwsLogsAPI struct {
	awsCfg aws.Config
}

func NewAwsLogsAPI(c aws.Config) *AwsLogsAPI {
	return &AwsLogsAPI{
		awsCfg: c,
	}
}

func (api AwsLogsAPI) Validate(resource string) bool {
	switch resource {
	case "log-group":
		return true
	}

	return false
}

func (api AwsLogsAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service: "logs",
		Resource: resource,
	}

	awsAPI := cloudwatchlogs.NewFromConfig(api.awsCfg)
	switch resource {
	case "log-group":
		listOutput, err := awsAPI.DescribeLogGroups(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		for _, lg := range listOutput.LogGroups {
			resultList.Results = append(resultList.Results, lg)
		}
	default:
		log.Error().Msgf("resource '%s' not supported in logs service", resource)
		return nil, fmt.Errorf("resource '%s' not supported in logs service", resource)
	}

	return resultList, nil
}
