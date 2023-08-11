package internal

import (
	"context"
	"fmt"
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog/log"
)

type AwsresqClient struct {
	awsCfg aws.Config
	Region []string
}

type ResultList struct {
	Service string `json:"service"`
	Resource string `json:"resource"`
	Results []interface{} `json:"results"`
}

func NewAwsresqClient(region string) (*AwsresqClient, error) {
	client := &AwsresqClient{}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration error")
		return nil, err
	}
	client.awsCfg = cfg

	client.Region = buildRegion(region)

	return client, nil
}

func (c *AwsresqClient) Search(service, resource string) (string, error) {
	var resultList *ResultList

	switch service {
	case "ecs":
		api := NewAwsEcsAPI(c.awsCfg, c.Region)
		var err error
		resultList, err = api.Query(resource)
		if err != nil {
			return "", err
		}
	case "logs":
		api := NewAwsLogsAPI(c.awsCfg)
		var err error
		resultList, err = api.Query(resource)
		if err != nil {
			return "", err
		}
	default:
		log.Error().Msgf("service not supported: %s", service)
		return "", fmt.Errorf("service not supported: %s", service)
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
