package manager

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/controller"
	"advanced-tools/pkg/entity"
	"fmt"
	"strings"
)

type KubernetesUpgradeManager struct {
	k8sUpgradeController        *controller.K8sUpgradeController
	slackNotificationController *controller.SlackNotificationController
}

func GetKubernetesUpgradeManager(clients *client.Client) *KubernetesUpgradeManager {
	return &KubernetesUpgradeManager{
		k8sUpgradeController:        controller.GetK8sUpgradeController(clients),
		slackNotificationController: controller.GetSlackNotificationController(clients),
	}
}

func (manager *KubernetesUpgradeManager) Run() {
	// pre process
	// lock loop on request to start upgrade process by a devops team member - input should contain "target version: [target version]"
	// post to channel "upgrade process is starting, stage [pre-process-validations], target version [provided target version]"
	asgs, err := manager.k8sUpgradeController.GetASGsNodeList()
	if err != nil {
		return
	}
	asgsStr := formatOutAsgNodeList(asgs)
	err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(asgsStr)
	if err != nil {
		return
	}
	// verify all nodes are in the same version
	// post result to channel, wait for response: abort or proceed
	// perform search of all crash loop / image pull err / errored pods in cluster
	// post result to channel, wait for response: abort or proceed
	// list all alerts that are going to be silenced, wait for approval to silence them
	// silence alerts, post to channel that alerts been silenced
	// post to dev channel with link to notifications channel
	// process
	// upgrade process continues to next stage [upgrading]
	// git pull tf project
	// modify relevant file
	// tf plan
	// post results to slack, wait for approval or abort
	// tf apply
	// git add + commit + push + provide link to the new branch, post on channel
	// post table of all asgs, notify "starting rollout restart"
	// rolling nodes of asg [asg name] [asg label]
	// every 5 min trigger check of num nodes that been upgraded and post it by percentage
	// post process
	// upgrade process is in finalize state
	// post report of all non aligned nodes if any
	// post report of imagepull / crashloop / errored pods in cluster
}

func formatOutAsgNodeList(asgList []*entity.ASGNodeList) string {
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
