//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE_mock -destination=../mock/$GOFILE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type awsIamAPI interface {
	ListRoles(ctx context.Context, params *iam.ListRolesInput, optFns ...func(*iam.Options)) (*iam.ListRolesOutput, error)
	ListUsers(ctx context.Context, params *iam.ListUsersInput, optFns ...func(*iam.Options)) (*iam.ListUsersOutput, error)
}

type AwsresqIamAPI struct {
	awsCfg    aws.Config
	region    []string
	apiClient map[string]awsIamAPI
}

func NewAwsresqIamAPI(c aws.Config, region []string) *AwsresqIamAPI {
	return &AwsresqIamAPI{
		awsCfg:    c,
		region:    region,
		apiClient: make(map[string]awsIamAPI, len(region)),
	}
}

func (api AwsresqIamAPI) Validate(resource string) bool {
	validResource := []string{
		"role",
		"user",
	}
	return slices.Contains(validResource, resource)
}

func (api AwsresqIamAPI) Query(resource string) (*ResultList, error) {
	resultList := &ResultList{
		Service:  "iam",
		Resource: resource,
	}

	var apiQuery ResourceQueryAPI
	switch resource {
	case "role":
		apiQuery = api.queryIamRole
		api.region = []string{"us-east-1"}
	case "user":
		apiQuery = api.queryIamUser
		api.region = []string{"us-east-1"}
	default:
		return nil, fmt.Errorf("resource %s is not supported in iam service", resource)
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

func (api AwsresqIamAPI) queryIamRole(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service:  "iam",
		Resource: "role",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = iam.NewFromConfig(api.awsCfg, func(o *iam.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].ListRoles(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to list roles in region %s", region)
		return
	}
	for _, role := range listOutput.Roles {
		// AssumeRolePolicyDocument is URL encoded. It needs to be unescaped as below for query result.
		// doc, _ := url.PathUnescape(*role.AssumeRolePolicyDocument)
		// role.AssumeRolePolicyDocument = aws.String(doc)
		resultList.Results = append(resultList.Results, role)
	}

	ch <- resultList
}

func (api AwsresqIamAPI) queryIamUser(ctx context.Context, ch chan ResultList, region string) {
	resultList := ResultList{
		Service:  "iam",
		Resource: "user",
	}

	if api.apiClient[region] == nil {
		api.apiClient[region] = iam.NewFromConfig(api.awsCfg, func(o *iam.Options) {
			o.Region = region
		})
	}

	listOutput, err := api.apiClient[region].ListUsers(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to list users in region %s", region)
		return
	}
	for _, user := range listOutput.Users {
		resultList.Results = append(resultList.Results, user)
	}

	ch <- resultList
}
