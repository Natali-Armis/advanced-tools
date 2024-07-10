package controller

import (
	"advanced-tools/pkg/entity"
	"fmt"
	"strings"
)

func FormatOutAsgNodeList(header string, asgList []*entity.ASGNodeList) string {
	var builder strings.Builder

	maxAsgName := len("AsgName")
	maxLabel := len("Label")
	maxInstanceId := len("InstanceId")
	maxPrivateDnsName := len("PrivateDnsName")
	maxKubeletVersion := len("KubeletVersion")

	for _, asg := range asgList {
		if len(asg.AsgName) > maxAsgName {
			maxAsgName = len(asg.AsgName)
		}
		if len(asg.Label) > maxLabel {
			maxLabel = len(asg.Label)
		}
		for _, node := range asg.NodeList {
			if len(node.InstanceId) > maxInstanceId {
				maxInstanceId = len(node.InstanceId)
			}
			if len(node.PrivateDnsName) > maxPrivateDnsName {
				maxPrivateDnsName = len(node.PrivateDnsName)
			}
			if len(node.KubeletVersion) > maxKubeletVersion {
				maxKubeletVersion = len(node.KubeletVersion)
			}
		}
	}

	fmt.Fprintf(&builder, "%v:\n", header)
	fmt.Fprintf(&builder, "%v\n", strings.Repeat("-", maxAsgName+maxLabel+maxInstanceId+maxPrivateDnsName+maxKubeletVersion+50))

	for _, asg := range asgList {
		for _, node := range asg.NodeList {
			fmt.Fprintf(&builder, "%-*s %-*s %-*s %-*s %-*s\n",
				maxAsgName, asg.AsgName,
				maxLabel, asg.Label,
				maxInstanceId, node.InstanceId,
				maxPrivateDnsName, node.PrivateDnsName,
				maxKubeletVersion, node.KubeletVersion)
		}
	}
	return "\n" + builder.String() + "\n"
}

func FormatFailingPodsList(header string, failingPods []*entity.FailingPod) string {
	var builder strings.Builder

	maxPodName := len("PodName")
	maxNamespace := len("Namespace")
	maxStatus := len("Status")

	for _, pod := range failingPods {
		if len(pod.PodName) > maxPodName {
			maxPodName = len(pod.PodName)
		}
		if len(pod.Namespace) > maxNamespace {
			maxNamespace = len(pod.Namespace)
		}
		if len(pod.Status) > maxStatus {
			maxStatus = len(pod.Status)
		}
	}

	fmt.Fprintf(&builder, "%v:\n", header)
	fmt.Fprintf(&builder, "%v\n", strings.Repeat("-", maxPodName+maxNamespace+maxStatus+20))

	for _, pod := range failingPods {
		fmt.Fprintf(&builder, "%-*s %-*s %-*s\n",
			maxPodName, pod.PodName,
			maxNamespace, pod.Namespace,
			maxStatus, pod.Status)
	}
	return "\n" + builder.String() + "\n"
}