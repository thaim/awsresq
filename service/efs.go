//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsEfsAPI interface {
	DescribeFileSystems(ctx context.Context, params *efs.DescribeFileSystemsInput, optFns ...func(*efs.Options)) (*efs.DescribeFileSystemsOutput, error)
}

type AwsresqEfsAPI struct {
	awsCfg aws.Config
    region []string
	apiClient map[string]awsEfsAPI
}

func NewAwsresqEfsAPI(awsCfg aws.Config, region []string) *AwsresqEfsAPI {
	return &AwsresqEfsAPI{
		awsCfg: awsCfg,
		region: region,
		apiClient: make(map[string]awsEfsAPI, len(region)),
	}
}

func (a *AwsresqEfsAPI) Validate(resource string) bool {
	validResources := []string{
		"file-system",
	}

	return slices.Contains(validResources, resource)
}

func (api AwsresqEfsAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service: "efs",
		Resource: resource,
	}
	var err error = nil

	var apiQuery ResourceQueryAPI
	switch resource {
	case "file-system":
		apiQuery = api.queryFileSystem
	default:
		log.Error().Msgf("resource %s not supported in efs service", resource)
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
		case result := <-ch:
			resultList.Results = append(resultList.Results, result.Results...)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return resultList, nil
}

func (api *AwsresqEfsAPI) queryFileSystem(ctx context.Context, ch chan ResultList, r string) {
	resultList := ResultList{
		Service: "efs",
		Resource: "file-system",
	}

	if api.apiClient[r] == nil {
		api.apiClient[r] = efs.NewFromConfig(api.awsCfg, func(o *efs.Options) {
			o.Region = r
		})
	}

	listOutput, err := api.apiClient[r].DescribeFileSystems(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to describe file systems in %s", r)
		return
	}
	resultList.Results = append(resultList.Results, listOutput.FileSystems)

	ch <- resultList
}
