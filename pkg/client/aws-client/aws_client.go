package aws_client

import (
	"context"
	"strings"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	as "github.com/aws/aws-sdk-go-v2/service/autoscaling"
	as_types "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/rs/zerolog/log"
)

type AwsClient struct {
	eksClient *eks.Client
	asgClient *as.Client
	ec2Client *ec2.Client
	region    string
}

func GetAwsClient(profile string, region string) *AwsClient {
	log.Info().Msgf("configuring aws client")
	log.Info().Msg("loading default AWS profile")
	awsCfg, err := aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithRegion(region))
	if len(profile) > 0 {
		log.Info().Msgf("loading AWS profile: %s", profile)
		awsCfg, err = aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithSharedConfigProfile(profile), aws_config.WithRegion(region))
	}

	if err != nil {
		log.Fatal().Msgf("unable to load SDK config, %v", err)
	}
	eksClient := eks.NewFromConfig(awsCfg)
	asgClient := as.NewFromConfig(awsCfg)
	ec2Client := ec2.NewFromConfig(awsCfg)

	log.Info().Msgf("aws client configured, profile: %v, region: %v", profile, region)
	return &AwsClient{
		eksClient: eksClient,
		asgClient: asgClient,
		ec2Client: ec2Client,
		region:    region,
	}
}

func (client *AwsClient) DescribeAutoScalingGroups(clusterNameSubstring string) ([]as_types.AutoScalingGroup, error) {
	output, err := client.asgClient.DescribeAutoScalingGroups(context.TODO(), &as.DescribeAutoScalingGroupsInput{})
	if err != nil {
		log.Error().Msgf("error during describing cluster %v autoscaling groups %v", clusterNameSubstring, err.Error())
		return nil, err
	}

	var asgs []as_types.AutoScalingGroup
	for _, asg := range output.AutoScalingGroups {
		if strings.Contains(*asg.AutoScalingGroupName, clusterNameSubstring) {
			asgs = append(asgs, asg)
		}
	}
	return asgs, nil
}

func (client *AwsClient) ListInstances(asgName string) ([]string, error) {
	output, err := client.asgClient.DescribeAutoScalingGroups(context.TODO(), &as.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil {
		log.Error().Msgf("error during listing asg %v instancess %v", asgName, err.Error())
		return nil, err
	}
	var instanceIDs []string
	for _, asg := range output.AutoScalingGroups {
		for _, instance := range asg.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}
	return instanceIDs, nil
}

func (client *AwsClient) DescribeInstances(instanceIDs []string) ([]ec2_types.Instance, error) {
	output, err := client.ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		log.Error().Msgf("error during describing instances ids %v", err.Error())
		return nil, err
	}

	var instances []ec2_types.Instance
	for _, reservation := range output.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances, nil
}
