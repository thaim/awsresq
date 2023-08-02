package internal

import (
	"context"
	"fmt"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog/log"
)

type AwsresqClient struct {
	awsCfg aws.Config
}

type ResultList struct {
	Service string `json:service`
	Resource string `json:resource`
	Query string `json:query`
	Results []interface{} `json:"results"`
}

func NewAwsresqClient() (*AwsresqClient, error) {
	client := &AwsresqClient{}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration error")
		return nil, err
	}
	client.awsCfg = cfg

	return client, nil
}

func (c *AwsresqClient) Search(service, resource, query string) (string, error) {
	var resultList *ResultList

	switch service {
	case "ecs":
		api := NewAwsEcsAPI(c.awsCfg)
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
