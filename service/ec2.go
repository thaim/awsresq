//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsEc2API interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
	DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
}

type AwsresqEc2API struct {
	awsCfg    aws.Config
	region    []string
	apiClient map[string]awsEc2API
}

func NewAwsresqEc2API(c aws.Config, region []string) *AwsresqEc2API {
	return &AwsresqEc2API{
		awsCfg:    c,
		region:    region,
		apiClient: make(map[string]awsEc2API, len(region)),
	}
}

func (api AwsresqEc2API) Validate(resource string) bool {
	validResource := []string{
		"instance",
		"security-group",
		"vpc",
	}
	return slices.Contains(validResource, resource)
}

func (api AwsresqEc2API) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "ec2",
		Resource: resource,
	}

	var apiQuery ResourceQueryAPI
	switch resource {
	case "instance":
		apiQuery = api.queryEc2Instance
	case "security-group":
		apiQuery = api.queryEc2SecurityGroup
	case "vpc":
		apiQuery = api.queryEc2Vpc
	default:
		return nil, fmt.Errorf("resource %s is not supported in ec2 service", resource)
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

func (api AwsresqEc2API) queryEc2Instance(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service:  "ec2",
		Resource: "instance",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = ec2.NewFromConfig(api.awsCfg, func(o *ec2.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].DescribeInstances(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to describe ec2 instance in region %s", region)
		return
	}
	for _, reservation := range listOutput.Reservations {
		for _, instance := range reservation.Instances {
			resultList.Results = append(resultList.Results, instance)
		}
	}

	ch <- resultList
}

func (api AwsresqEc2API) queryEc2SecurityGroup(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service:  "ec2",
		Resource: "security-group",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = ec2.NewFromConfig(api.awsCfg, func(o *ec2.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].DescribeSecurityGroups(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to describe ec2 security group in region %s", region)
		return
	}
	for _, securityGroup := range listOutput.SecurityGroups {
		resultList.Results = append(resultList.Results, securityGroup)
	}

	ch <- resultList
}

func (api AwsresqEc2API) queryEc2Vpc(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service:  "ec2",
		Resource: "vpc",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = ec2.NewFromConfig(api.awsCfg, func(o *ec2.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].DescribeVpcs(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to describe ec2 vpc in region %s", region)
		return
	}
	for _, vpc := range listOutput.Vpcs {
		resultList.Results = append(resultList.Results, vpc)
	}

	ch <- resultList
}
