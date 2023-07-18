package internal

import (
	"context"
	"fmt"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/rs/zerolog/log"
)

type AwsresqClient struct {
	awsCfg aws.Config
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

func (c *AwsresqClient) Search(service, resource, query string) ([]string, error) {
	var resultList []string
	var result []byte

	switch service {
	case "ecs":
		api := ecs.NewFromConfig(c.awsCfg)
		switch resource {
		case "task-definition":
			listOutput, err := api.ListTaskDefinitions(context.Background(), nil)
			if err != nil {
				return resultList, err
			}
			for _, arn := range listOutput.TaskDefinitionArns {
				input := &ecs.DescribeTaskDefinitionInput{
					TaskDefinition: aws.String(arn),
				}
				output, err := api.DescribeTaskDefinition(context.Background(), input)
				if err != nil {
					return resultList, err
				}

				result, err = json.Marshal(output)
				if err != nil {
					return resultList, err
				}

				resultList = append(resultList, string(result))
			}
		default:
			log.Error().Msgf("resource '%s' not supported in service '%s'", resource, service)
			return resultList, fmt.Errorf("resource '%s' not supported in service '%s'", resource, service)
		}
	default:
		log.Error().Msgf("service not supported: %s", service)
		return resultList, fmt.Errorf("service not supported: %s", service)
	}


	return resultList, nil
}
