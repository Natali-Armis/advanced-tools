package aws_client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"advanced-tools/pkg/vars"

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
	log.Info().Msgf("client: configuring aws client")
	log.Info().Msg("client: loading default AWS profile")
	awsCfg, err := aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithRegion(region))
	if len(profile) > 0 {
		log.Info().Msgf("client: loading AWS profile [%s]", profile)
		awsCfg, err = aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithSharedConfigProfile(profile), aws_config.WithRegion(region))
	}

	if err != nil {
		log.Fatal().Msgf("client: unable to load SDK config, %v", err)
	}
	eksClient := eks.NewFromConfig(awsCfg)
	asgClient := as.NewFromConfig(awsCfg)
	ec2Client := ec2.NewFromConfig(awsCfg)

	log.Info().Msgf("client: aws client configured, profile [%v], region [%v]", profile, region)
	return &AwsClient{
		eksClient: eksClient,
		asgClient: asgClient,
		ec2Client: ec2Client,
		region:    region,
	}
}

func (client *AwsClient) DescribeAutoScalingGroups(clusterNameSubstring string) ([]as_types.AutoScalingGroup, error) {
	clusterNameSubstring = strings.ReplaceAll(clusterNameSubstring, "-", "_")
	clusterNameSubstring = fmt.Sprintf("%v_%v", clusterNameSubstring, vars.ASG_REQUIRED_LABEL)
	output, err := client.asgClient.DescribeAutoScalingGroups(context.TODO(), &as.DescribeAutoScalingGroupsInput{})
	if err != nil {
		log.Error().Msgf("client: error during describing cluster %v autoscaling groups %v", clusterNameSubstring, err.Error())
		return nil, err
	}

	asgs := []as_types.AutoScalingGroup{}
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
		log.Error().Msgf("client: error during listing asg %v instancess %v", asgName, err.Error())
		return nil, err
	}
	instanceIDs := []string{}
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
		log.Error().Msgf("client: error during describing instances ids %v", err.Error())
		return nil, err
	}

	instances := []ec2_types.Instance{}
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			if _, exists := vars.AwsInstancesCodes[*instance.State.Code]; exists {
				instances = append(instances, instance)
			}
		}
	}
	return instances, nil
}

func (client *AwsClient) GetMaxSizeOfASG(asgName string) (int32, error) {
	output, err := client.asgClient.DescribeAutoScalingGroups(context.TODO(), &as.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil {
		log.Error().Msgf("client: error during getting max size of asg %v, %v", asgName, err.Error())
		return 0, err
	}
	if len(output.AutoScalingGroups) == 0 {
		return 0, fmt.Errorf("client: no ASG found with name %v", asgName)
	}
	return *output.AutoScalingGroups[0].MaxSize, nil
}

func (client *AwsClient) ModifyMaxSizeOfASG(asgName string, newSize int32) error {
	_, err := client.asgClient.UpdateAutoScalingGroup(context.TODO(), &as.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: &asgName,
		MaxSize:              &newSize,
	})
	if err != nil {
		log.Error().Msgf("client: error during modifying max size of asg %v, %v", asgName, err.Error())
		return err
	}
	return nil
}

func (client *AwsClient) GetDesiredSizeOfASG(asgName string) (int32, error) {
	output, err := client.asgClient.DescribeAutoScalingGroups(context.TODO(), &as.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil {
		log.Error().Msgf("client: error during getting desired size of asg %v, %v", asgName, err.Error())
		return 0, err
	}
	if len(output.AutoScalingGroups) == 0 {
		return 0, fmt.Errorf("client: no ASG found with name %v", asgName)
	}
	return *output.AutoScalingGroups[0].DesiredCapacity, nil
}

func (client *AwsClient) ModifyDesiredSizeOfASG(asgName string, newSize int32) error {
	_, err := client.asgClient.UpdateAutoScalingGroup(context.TODO(), &as.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: &asgName,
		DesiredCapacity:      &newSize,
	})
	if err != nil {
		log.Error().Msgf("client: error during modifying desired size of asg %v, %v", asgName, err.Error())
		return err
	}
	return nil
}

func (client *AwsClient) GetInstancesWithScaleInProtection(asgName string) ([]string, error) {
	output, err := client.asgClient.DescribeAutoScalingGroups(context.TODO(), &as.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil {
		log.Error().Msgf("client: error during getting instances with scale-in protection for asg %v, %v", asgName, err.Error())
		return nil, err
	}

	instanceIDs := []string{}
	for _, asg := range output.AutoScalingGroups {
		for _, instance := range asg.Instances {
			if *instance.ProtectedFromScaleIn {
				instanceIDs = append(instanceIDs, *instance.InstanceId)
			}
		}
	}
	return instanceIDs, nil
}

func (client *AwsClient) RemoveScaleInProtection(asgName string, instanceID string) error {
	protected := false
	_, err := client.asgClient.SetInstanceProtection(context.TODO(), &as.SetInstanceProtectionInput{
		AutoScalingGroupName: &asgName,
		InstanceIds:          []string{instanceID},
		ProtectedFromScaleIn: &protected,
	})
	if err != nil {
		log.Error().Msgf("client: error during removing scale-in protection for instance %v in asg %v, %v", instanceID, asgName, err.Error())
		return err
	}
	return nil
}

func (client *AwsClient) AddScaleInProtection(asgName string, instanceID string) error {
	protected := true
	_, err := client.asgClient.SetInstanceProtection(context.TODO(), &as.SetInstanceProtectionInput{
		AutoScalingGroupName: &asgName,
		InstanceIds:          []string{instanceID},
		ProtectedFromScaleIn: &protected,
	})
	if err != nil {
		log.Error().Msgf("client: error during adding scale-in protection for instance %v in asg %v, %v", instanceID, asgName, err.Error())
		return err
	}
	return nil
}

func (client *AwsClient) GetNewInstances(asgName string, createdAfter time.Time) ([]string, error) {
	output, err := client.asgClient.DescribeAutoScalingGroups(context.TODO(), &as.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil {
		log.Error().Msgf("client: error during getting new instances for asg %v, %v", asgName, err.Error())
		return nil, err
	}
	instanceIDs := []string{}
	for _, asg := range output.AutoScalingGroups {
		for _, instance := range asg.Instances {
			instanceOutput, err := client.ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
				InstanceIds: []string{*instance.InstanceId},
			})
			if err != nil {
				log.Error().Msgf("client: error during describing instance %v, %v", *instance.InstanceId, err.Error())
				return nil, err
			}
			for _, reservation := range instanceOutput.Reservations {
				for _, instanceDetail := range reservation.Instances {
					if instanceDetail.LaunchTime.After(createdAfter) {
						instanceIDs = append(instanceIDs, *instance.InstanceId)
					}
				}
			}
		}
	}
	return instanceIDs, nil
}
