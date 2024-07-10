package manager

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/controller"
	"advanced-tools/pkg/vars"
	"fmt"
	"regexp"
	"time"
)

type KubernetesUpgradeManager struct {
	k8sUpgradeController        *controller.K8sUpgradeController
	slackNotificationController *controller.SlackNotificationController
	alertsController            *controller.AlertsController
	stage                       string
}

func GetKubernetesUpgradeManager(clients *client.Client) *KubernetesUpgradeManager {
	return &KubernetesUpgradeManager{
		k8sUpgradeController:        controller.GetK8sUpgradeController(clients),
		slackNotificationController: controller.GetSlackNotificationController(clients),
		alertsController:            controller.GetAlertsController(clients),
		stage:                       "",
	}
}

func (manager *KubernetesUpgradeManager) Run() {
	// pre process
	targetVersion := ""
	for len(targetVersion) == 0 {
		lastMessageFromChannel, userId, err := manager.slackNotificationController.FetchLastMessageFromUpgradeNotificationsChannelMatchPattern(vars.UPGEADE_INITIALIZE_SUBSTR)
		if err != nil {
			return
		}
		if len(lastMessageFromChannel) > 0 {
			if !manager.slackNotificationController.IdentifySelfMessage(userId) {
				targetVersion = extractTargetVersion(lastMessageFromChannel)
			}
			if len(targetVersion) > 0 {
				break
			}
		} else {
			err := manager.slackNotificationController.NotifyInUpgradeNotificationsChannel("kubernetes manager for upgrade service, usage: `kubernetes-manager upgrade cluster <target version>`\n(i.e target version must be formatted as `\\d+\\.\\d+` example `1.30`)")
			if err != nil {
				return
			}
		}
		time.Sleep(10 * time.Second)
	}

	err := manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(fmt.Sprintf("upgrade process is starting, stage [%v], target version [%v]", vars.PRE_PROCESS, targetVersion))
	if err != nil {
		return
	}

	asgs, err := manager.k8sUpgradeController.GetASGsNodeList()
	if err != nil {
		return
	}
	asgsStr := controller.FormatOutAsgNodeList("listing cluster autoscaling groups and their owned nodes", asgs)
	for _, str := range asgsStr {
		err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(str)
		if err != nil {
			return
		}
	}

	// err = manager.alertsController.SilenceAlerts()
	// if err != nil {
	// 	return
	// }
	failingPods, err := manager.k8sUpgradeController.GetErroredPodsList()
	if err != nil {
		return
	}
	failingPodsStr := controller.FormatFailingPodsList("listing failing pods in cluster", failingPods)
	for _, str := range failingPodsStr {
		err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(str)
		if err != nil {
			return
		}
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

func extractTargetVersion(message string) string {
	re := regexp.MustCompile(vars.UPGRADE_INITIALIZE_PATTERN)
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
