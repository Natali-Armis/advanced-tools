package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"
	"advanced-tools/pkg/entity"
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
	asgs, err := getASGsNodeList()
	if err != nil {
		return
	}
	printOutAsgNodeList(asgs)
}

func printOutAsgNodeList(asgList []*entity.ASGNodeList) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(writer, "AutoScaling Group\t Service Label\t Instance ID\t Private DNS\t Kubelet Version")
	fmt.Fprintln(writer, "-----------------\t -------------\t -----------\t ------------\t --------------")
	for _, asg := range asgList {
		for _, node := range asg.NodeList {
			fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\n", asg.AsgName, asg.Label, node.InstanceId, node.PrivateDnsName, node.KubeletVersion)
		}
	}
	writer.Flush()
}

func getASGsNodeList() ([]*entity.ASGNodeList, error) {
	asgList := []*entity.ASGNodeList{}
	clusterAsgs, err := clients.AwsClient.DescribeAutoScalingGroups(vars.DEV)
	if err != nil {
		return asgList, err
	}
	for _, clusterAsg := range clusterAsgs {
		asg := &entity.ASGNodeList{
			AsgName:  *clusterAsg.AutoScalingGroupName,
			NodeList: []*entity.AsgNode{},
		}
		asgInstances, err := clients.AwsClient.ListInstances(*clusterAsg.AutoScalingGroupName)
		if err != nil {
			return asgList, err
		}
		instancesDescribe, err := clients.AwsClient.DescribeInstances(asgInstances)
		if err != nil {
			return asgList, err
		}
		for _, instance := range instancesDescribe {
			version, label, err := clients.K8sClient.GetNodeVersionAndLabels(*instance.PrivateDnsName)
			if err != nil {
				continue
			}
			if len(asg.Label) == 0 {
				asg.Label = label
			}
			asg.NodeList = append(asg.NodeList, &entity.AsgNode{
				InstanceId:     *instance.InstanceId,
				PrivateDnsName: *instance.PrivateDnsName,
				KubeletVersion: version,
			})
		}
		asgList = append(asgList, asg)
	}
	return asgList, nil
}
