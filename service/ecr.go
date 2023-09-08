package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsEcrAPI interface {
	DescribeRepositories(ctx context.Context, params *ecr.DescribeRepositoriesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error)
}

type AwsresqEcrAPI struct {
	awsCfg    aws.Config
	region    []string
	apiClient map[string]awsEcrAPI
}

func NewAwsresqEcrAPI(c aws.Config, region []string) AwsresqEcrAPI {
	return AwsresqEcrAPI{
		awsCfg:    c,
		region:    region,
		apiClient: make(map[string]awsEcrAPI, len(region)),
	}
}

func (api AwsresqEcrAPI) Validate(resource string) bool {
	validResources := []string{
		"repository",
	}

	return slices.Contains(validResources, resource)
}

func (api AwsresqEcrAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "ecr",
		Resource: resource,
	}
	var err error = nil

	var apiQuery ResourceQueryAPI
	switch resource {
	case "repository":
		apiQuery = api.queryRepository
	default:
		log.Error().Msgf("resource %s not supported in ecr service", resource)
		return nil, err
	}

	ch := make(chan ResultList)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, r := range api.region {
		go apiQuery(ctx, ch, r)
	}

	for range api.region {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("ecr: %v", ctx.Err())
			log.Error().Msgf("ecr: %v", ctx.Err())
		case res := <-ch:
			resultList.Results = append(resultList.Results, res.Results...)
		}
	}

	return resultList, err
}

func (api AwsresqEcrAPI) queryRepository(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service:  "ecr",
		Resource: "repository",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = ecr.NewFromConfig(api.awsCfg, func(o *ecr.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].DescribeRepositories(ctx, nil)
	if err != nil {
		log.Error().Msgf("error querying ecr repository in %s: %v", region, err)
		return
	}
	for _, repo := range listOutput.Repositories {
		resultList.Results = append(resultList.Results, repo)
	}

	ch <- resultList
}
