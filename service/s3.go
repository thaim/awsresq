//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsS3API interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

type AwsresqS3API struct {
	awsCfg aws.Config
	region []string
	apiClient map[string]awsS3API
}

func NewAwsresqS3API(awsConfig aws.Config, region []string) *AwsresqS3API {
	return &AwsresqS3API{
		awsCfg: awsConfig,
		region: region,
		apiClient: make(map[string]awsS3API, len(region)),
	}
}

func (api AwsresqS3API) Validate(resource string) bool {
	validResource := []string{
		"bucket",
	}

	return slices.Contains(validResource, resource)
}

func (api AwsresqS3API) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service: "s3",
		Resource: resource,
	}
	var apiQuery ResourceQueryAPI

	switch resource {
	case "bucket":
		apiQuery = api.queryBucket
	default:
		return nil, fmt.Errorf("resource %s not supported in s3 service", resource)
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

func (api *AwsresqS3API) queryBucket(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service: "s3",
		Resource: "bucket",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = s3.NewFromConfig(api.awsCfg, func(o *s3.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].ListBuckets(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to list bucket in %s", region)
		return
	}
	if len(listOutput.Buckets) > 0 {
		for _, b := range listOutput.Buckets {
			resultList.Results = append(resultList.Results, b)
		}
	}

	ch <- resultList
}
