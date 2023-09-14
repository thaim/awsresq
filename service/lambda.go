//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsLambdaAPI interface {
	ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

type AwsresqLambdaAPI struct {
	awsCfg    aws.Config
	region    []string
	apiClient map[string]awsLambdaAPI
}

func NewAwsresqLambdaAPI(c aws.Config, region []string) *AwsresqLambdaAPI {
	return &AwsresqLambdaAPI{
		awsCfg:    c,
		region:    region,
		apiClient: make(map[string]awsLambdaAPI, len(region)),
	}
}

func (api AwsresqLambdaAPI) Validate(resource string) bool {
	validResources := []string{
		"function",
	}

	return slices.Contains(validResources, resource)
}

func (api AwsresqLambdaAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "lambda",
		Resource: resource,
	}

	var apiQuery ResourceQueryAPI
	switch resource {
	case "function":
		apiQuery = api.queryFunction
	default:
		log.Error().Msgf("resource '%s' not supported in lambda service", resource)
		return nil, fmt.Errorf("resource '%s' not supported in lambda service", resource)
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
			if result.Results != nil {
				resultList.Results = append(resultList.Results, result.Results...)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return resultList, nil
}

func (api *AwsresqLambdaAPI) queryFunction(ctx context.Context, ch chan ResultList, r string) {
	resultList := ResultList{
		Service:  "lambda",
		Resource: "function",
	}

	if api.apiClient[r] == nil {
		api.apiClient[r] = lambda.NewFromConfig(api.awsCfg, func(o *lambda.Options) {
			o.Region = r
		})
	}

	listOutput, err := api.apiClient[r].ListFunctions(ctx, nil)
	if err != nil {
		log.Error().Msgf("failed to list functions in %s: %s", r, err.Error())
		return
	}
	for _, function := range listOutput.Functions {
		resultList.Results = append(resultList.Results, function)
	}

	ch <- resultList
}
