package manager

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/controller"
	"advanced-tools/pkg/vars"
	"fmt"
	"regexp"
	"time"

	"github.com/rs/zerolog/log"
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
	targetVersion := manager.waitForUserToStart()

	err := manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(fmt.Sprintf("upgrade process is starting, stage [%v], target version [%v]", vars.PRE_PROCESS, targetVersion))
	if err != nil {
		return
	}

	proceed := manager.managePreProcessValidations()
	if !proceed {
		log.Warn().Msgf("aborting process...")
		return
	}

	err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(fmt.Sprintf("upgrade process is running, stage [%v], target version [%v]", vars.EDITING_TF_FILES, targetVersion))
	if err != nil {
		return
	}

	// git pull tf project
	// modify relevant file
	// tf plan
	// post results to slack, wait for approval or abort

	err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(fmt.Sprintf("upgrade process is running, stage [%v], target version [%v]", vars.PENDING_APPROVAL, targetVersion))
	if err != nil {
		return
	}

	// receive approval or abort

	err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(fmt.Sprintf("upgrade process is running, stage [%v], target version [%v]", vars.APPLYING_TF, targetVersion))
	if err != nil {
		return
	}

	// tf apply
	// git add + commit + push + provide link to the new branch, post on channel

	err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(fmt.Sprintf("upgrade process is running, stage [%v], target version [%v]", vars.ROLLING_EKS_NODES, targetVersion))
	if err != nil {
		return
	}

	// post table of all asgs, notify "starting rollout restart"
	// rolling nodes of asg [asg name] [asg label]
	// every 5 min trigger check of num nodes that been upgraded and post it by percentage

	err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(fmt.Sprintf("upgrade process is running, stage [%v], target version [%v]", vars.POST_PROCESS, targetVersion))
	if err != nil {
		return
	}

	// upgrade process is in finalize state
	// post report of all non aligned nodes if any
	// post report of imagepull / crashloop / errored pods in cluster
}

func (manager *KubernetesUpgradeManager) waitForUserToStart() string {
	targetVersion := ""
	for len(targetVersion) == 0 {
		lastMessageFromChannel, userId, err := manager.slackNotificationController.FetchLastMessageFromUpgradeNotificationsChannelMatchPattern(vars.UPGEADE_INITIALIZE_SUBSTR)
		if err != nil {
			return ""
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
				return ""
			}
		}
		time.Sleep(10 * time.Second)
	}
	return targetVersion
}

func (manager *KubernetesUpgradeManager) managePreProcessValidations() bool {
	// verify all nodes are in the same version
	asgs, err := manager.k8sUpgradeController.GetASGsNodeList()
	if err != nil {
		return false
	}
	asgsStr := controller.FormatOutAsgNodeList("listing cluster autoscaling groups and their owned nodes", asgs)
	for _, str := range asgsStr {
		err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(str)
		if err != nil {
			return false
		}
	}
	// post result to channel, wait for response: abort or proceed

	// perform search of all crash loop / image pull err / errored pods in cluster
	failingPods, err := manager.k8sUpgradeController.GetErroredPodsList()
	if err != nil {
		return false
	}
	failingPodsStr := controller.FormatFailingPodsList("listing failing pods in cluster", failingPods)
	for _, str := range failingPodsStr {
		err = manager.slackNotificationController.NotifyInUpgradeNotificationsChannel(str)
		if err != nil {
			return false
		}
	}
	// post result to channel, wait for response: abort or proceed

	// list all alerts that are going to be silenced, wait for approval to silence them
	// err = manager.alertsController.SilenceAlerts()
	// if err != nil {
	// 	return
	// }
	// silence alerts, post to channel that alerts been silenced
	// post to dev channel with link to notifications channel
	return true
}

func extractTargetVersion(message string) string {
	re := regexp.MustCompile(vars.UPGRADE_INITIALIZE_PATTERN)
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
