package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"
	"advanced-tools/pkg/vars"
	"fmt"
	"os"
	"text/tabwriter"
)

var (
	clients *client.Client
)

func init() {
	config.Configure()
	clients = client.GetClient()
}

func main() {

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(writer, "AutoScaling Group\tInstance ID\tPrivate DNS\tNode Name\tKubelet Version\tService Label")

	clusterAsgs, err := clients.AwsClient.DescribeAutoScalingGroups(vars.DEV)
	if err != nil {
		return
	}
	for _, clusterAsg := range clusterAsgs {
		asgInstances, err := clients.AwsClient.ListInstances(*clusterAsg.AutoScalingGroupName)
		if err != nil {
			return
		}
		instancesDescribe, err := clients.AwsClient.DescribeInstances(asgInstances)
		if err != nil {
			return
		}
		for _, instance := range instancesDescribe {
			version, label, err := clients.K8sClient.GetNodeVersionAndLabels(*instance.PrivateDnsName)
			if err != nil {
				continue
			}
			fmt.Printf("asg %v instance %v laebl %v version %v\n", *clusterAsg.AutoScalingGroupName, *instance.PrivateDnsName, label, version)
		}
	}

	
}
