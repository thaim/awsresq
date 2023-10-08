//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsRoute53API interface {
	ListHostedZones(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error)
}

type AwsresqRoute53API struct {
	awsCfg    aws.Config
	region    []string
	apiClient map[string]awsRoute53API
}

func NewAwsresqRoute53API(c aws.Config, region []string) *AwsresqRoute53API {
	return &AwsresqRoute53API{
		awsCfg:    c,
		region:    region,
		apiClient: make(map[string]awsRoute53API, len(region)),
	}
}

func (api AwsresqRoute53API) Validate(resource string) bool {
	validResoruces := []string{
		"hosted-zone",
	}

	return slices.Contains(validResoruces, resource)
}

func (api AwsresqRoute53API) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "route53",
		Resource: resource,
	}

	var apiQuery ResourceQueryAPI
	switch resource {
	case "hosted-zone":
		apiQuery = api.queryRoute53HostedZone
	default:
		return nil, fmt.Errorf("resource %s is not supported in ec2 service", resource)
	}

	ch := make(chan ResultList)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func (api AwsresqRoute53API) queryRoute53HostedZone(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service:  "route53",
		Resource: "hosted-zone",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = route53.NewFromConfig(api.awsCfg, func(o *route53.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].ListHostedZones(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("error listing hosted zones in %s", region)
		return
	}
	for _, hostedZone := range listOutput.HostedZones {
		resultList.Results = append(resultList.Results, hostedZone)
	}

	ch <- resultList
}
