//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsConfigAPI interface {
	DescribeConfigRules(ctx context.Context, params *configservice.DescribeConfigRulesInput, optFns ...func(*configservice.Options)) (*configservice.DescribeConfigRulesOutput, error)
}

type AwsresqConfigAPI struct {
	awsCfg aws.Config
	region []string
	apiClient map[string]awsConfigAPI
}

func NewAwsresqConfigAPI(awsCfg aws.Config, region []string) *AwsresqConfigAPI {
	return &AwsresqConfigAPI{
		awsCfg: awsCfg,
		region: region,
		apiClient: make(map[string]awsConfigAPI, len(region)),
	}
}

func (api AwsresqConfigAPI) Validate(resource string) bool {
	validResource := []string{
		"rule",
	}

	return slices.Contains(validResource, resource)
}

func (api AwsresqConfigAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service: "config",
		Resource: resource,
	}

	var apiQuery ResourceQueryAPI
	switch resource {
	case "rule":
		apiQuery = api.queryConfigRule
	default:
		return nil, fmt.Errorf("invalid resource type: %s", resource)
	}

	ch := make(chan ResultList)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, region := range api.region {
		go apiQuery(ctx, ch, region)
	}

	for range api.region {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-ch:
			resultList.Results = append(resultList.Results, result.Results...)
		}
	}

	return resultList, nil
}

func (api AwsresqConfigAPI) queryConfigRule(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service: "config",
		Resource: "rule",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = configservice.NewFromConfig(api.awsCfg, func(o *configservice.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].DescribeConfigRules(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to describe config rules")
		return
	}
	for _, rule := range listOutput.ConfigRules {
		resultList.Results = append(resultList.Results, rule)
	}

	ch <- resultList
}
