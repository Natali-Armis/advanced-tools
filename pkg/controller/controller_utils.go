package controller

import (
	"advanced-tools/pkg/entity"
	"fmt"
	"strings"
)

func FormatOutAsgNodeList(asgList []*entity.ASGNodeList) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%-140s %-40s %-30s %-60s %-20s\n", "AutoScaling Group", "Service Label", "Instance ID", "Private DNS", "Kubelet Version")
	fmt.Fprintf(&builder, "%-100s %-30s %-30s %-60s %-20s\n", strings.Repeat("-", 100), strings.Repeat("-", 30), strings.Repeat("-", 30), strings.Repeat("-", 60), strings.Repeat("-", 20))

	for _, asg := range asgList {
		for _, node := range asg.NodeList {
			fmt.Fprintf(&builder, "%-100s %-30s %-30s %-60s %-20s\n", asg.AsgName, asg.Label, node.InstanceId, node.PrivateDnsName, node.KubeletVersion)
		}
	}
	return "\n" + builder.String() + "\n"
}
