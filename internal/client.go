package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog/log"

	svc "github.com/thaim/awsresq/service"
)

type AwsresqClient struct {
	awsCfg aws.Config
	Region []string
	api    svc.AwsresqAPI
}

func NewAwsresqClient(region, service string) (*AwsresqClient, error) {
	client := &AwsresqClient{}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration error")
		return nil, err
	}
	client.awsCfg = cfg

	client.Region = buildRegion(region)

	switch service {
	case "cloudformation":
		client.api = svc.NewAwsresqCloudformationAPI(client.awsCfg, client.Region)
	case "cloudwatch":
		client.api = svc.NewAwsresqCloudwatchAPI(client.awsCfg, client.Region)
	case "config":
		client.api = svc.NewAwsresqConfigAPI(client.awsCfg, client.Region)
	case "ec2":
		client.api = svc.NewAwsresqEc2API(client.awsCfg, client.Region)
	case "ecr":
		client.api = svc.NewAwsresqEcrAPI(client.awsCfg, client.Region)
	case "ecs":
		client.api = svc.NewAwsresqEcsAPI(client.awsCfg, client.Region)
	case "efs":
		client.api = svc.NewAwsresqEfsAPI(client.awsCfg, client.Region)
	case "iam":
		client.api = svc.NewAwsresqIamAPI(client.awsCfg, client.Region)
	case "logs":
		client.api = svc.NewAwsresqLogsAPI(client.awsCfg, client.Region)
	case "lambda":
		client.api = svc.NewAwsresqLambdaAPI(client.awsCfg, client.Region)
	case "route53":
		client.api = svc.NewAwsresqRoute53API(client.awsCfg, client.Region)
	case "s3":
		client.api = svc.NewAwsresqS3API(client.awsCfg, client.Region)
	default:
		log.Error().Msgf("service not supported: %s", service)
		return nil, fmt.Errorf("service not supported: %s", service)
	}

	return client, nil
}

// Validate check if the resource is supported in the service
func (c *AwsresqClient) Validate(resource string) bool {
	return c.api.Validate(resource)
}

func (c *AwsresqClient) Search(service, resource string) (string, error) {
	var resultList *svc.ResultList
	resultList, err := c.api.Query(resource)
	if err != nil {
		return "", err
	}

	res, err := json.MarshalIndent(resultList, "", "  ")
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func buildRegion(region string) []string {
	if region == "all" || region == "" {
		return []string{
			"us-east-1", "us-east-2", "us-west-1", "us-west-2",
			"ap-south-1", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-southeast-1", "ap-southeast-2",
			"ca-central-1",
			"eu-central-1", "eu-west-1", "eu-west-2", "eu-west-3", "eu-north-1",
			"sa-east-1",
		}
	}

	return strings.Split(region, ",")
}
