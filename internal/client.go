package internal

import (
	"context"
	"fmt"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/itchyny/gojq"
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
	resultList := ResultList{
		Service: service,
		Resource: resource,
		Query: query,
	}

	goQuery, err := gojq.Parse(query)
	if err != nil {
		return "", err
	}

	switch service {
	case "ecs":
		api := ecs.NewFromConfig(c.awsCfg)
		switch resource {
		case "task-definition":
			listOutput, err := api.ListTaskDefinitions(context.Background(), nil)
			if err != nil {
				return "", err
			}
			for _, arn := range listOutput.TaskDefinitionArns {
				input := &ecs.DescribeTaskDefinitionInput{
					TaskDefinition: aws.String(arn),
					Include: []types.TaskDefinitionField{
						types.TaskDefinitionFieldTags,
					},
				}
				output, err := api.DescribeTaskDefinition(context.Background(), input)
				if err != nil {
					return "", err
				}

				var result interface{} = output
				iter := goQuery.Run(result)
				for {
					v, ok := iter.Next()
					if ! ok {
						break
					}
					if err, ok := v.(error); ok {
						log.Error().Err(err).Msg("failed to apply query")
						return "", err
					}
					resultList.Results = append(resultList.Results, v)
				}
			}
		default:
			log.Error().Msgf("resource '%s' not supported in service '%s'", resource, service)
			return "", fmt.Errorf("resource '%s' not supported in service '%s'", resource, service)
		}
	case "logs":
		api := cloudwatchlogs.NewFromConfig(c.awsCfg)
		switch resource {
		case "log-group":
			listOutput, err := api.DescribeLogGroups(context.Background(), nil)
			if err != nil {
				return "", err
			}
			for _, lg := range listOutput.LogGroups {
				resultList.Results = append(resultList.Results, lg)
			}
		default:
			log.Error().Msgf("resource '%s' not supported in service '%s'", resource, service)
			return "", fmt.Errorf("resource '%s' not supported in service '%s'", resource, service)
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
