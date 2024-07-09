package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"
	"advanced-tools/pkg/entity"
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

}

func printOutAsgNodeList(asgList []*entity.ASGNodeList) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(writer, "AutoScaling Group\t Service Label\t Instance ID\t Private DNS\t Kubelet Version")
	fmt.Fprintln(writer, "-----------------\t -------------\t -----------\t -----------\t ---------------")
	for _, asg := range asgList {
		for _, node := range asg.NodeList {
			fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\n", asg.AsgName, asg.Label, node.InstanceId, node.PrivateDnsName, node.KubeletVersion)
		}
	}
	writer.Flush()
}
