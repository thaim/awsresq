//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsCloudwatchAPI interface {
	ListMetrics(ctx context.Context, params *cloudwatch.ListMetricsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.ListMetricsOutput, error)
}

type AwsresqCloudwatchAPI struct {
	awsCfg    aws.Config
	region    []string
	apiClient map[string]awsCloudwatchAPI
}

func NewAwsresqCloudwatchAPI(c aws.Config, region []string) *AwsresqCloudwatchAPI {
	return &AwsresqCloudwatchAPI{
		awsCfg:    c,
		region:    region,
		apiClient: make(map[string]awsCloudwatchAPI, len(region)),
	}
}

func (api AwsresqCloudwatchAPI) Validate(resource string) bool {
	validResource := []string{
		"metric",
	}
	return slices.Contains(validResource, resource)
}

func (api AwsresqCloudwatchAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "cloudwatch",
		Resource: resource,
	}

	var apiQuery ResourceQueryAPI
	switch resource {
	case "metric":
		apiQuery = api.queryCloudwatchMetric
	default:
		return nil, fmt.Errorf("resource %s is not supported in cloudwatch service", resource)
	}

	ch := make(chan ResultList)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, region := range api.region {
		go apiQuery(ctx, ch, region)
	}

	for range api.region {
		select {
		case result := <-ch:
			resultList.Results = append(resultList.Results, result.Results...)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return resultList, nil
}

func (api *AwsresqCloudwatchAPI) queryCloudwatchMetric(ctx context.Context, ch chan ResultList, r string) {
	resultList := ResultList{
		Service:  "cloudwatch",
		Resource: "metric",
	}

	if api.apiClient[r] == nil {
		api.apiClient[r] = cloudwatch.NewFromConfig(api.awsCfg, func(o *cloudwatch.Options) {
			o.Region = r
		})
	}

	listOutput, err := api.apiClient[r].ListMetrics(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to list metrics in region %s", r)
		return
	}

	for _, metric := range listOutput.Metrics {
		resultList.Results = append(resultList.Results, metric)
	}

	ch <- resultList
}
