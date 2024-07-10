package controller

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/vars"
)

type SlackNotificationController struct {
	clients *client.Client
}

func GetSlackNotificationController(clients *client.Client) *SlackNotificationController {
	return &SlackNotificationController{
		clients: clients,
	}
}

func (controller *SlackNotificationController) NotifyInUpgradeNotificationsChannel(message string) error {
	return controller.clients.SlackClient.SendMessage(vars.SlackUpgradeNotificationsChannel, message)
}

func (controller *SlackNotificationController) FetchLastMessageFromUpgradeNotificationsChannel() (string, string, error) {
	return controller.clients.SlackClient.GetLastMessage(vars.SlackUpgradeNotificationsChannel)
}

func (controller *SlackNotificationController) FetchLastMessageFromUpgradeNotificationsChannelMatchPattern(pattern string) (string, string, error) {
	return controller.clients.SlackClient.GetLastMessageMatchPattern(vars.SlackUpgradeNotificationsChannel, pattern)
}

func (controller *SlackNotificationController) IdentifySelfMessage(userId string) bool {
	return userId == controller.clients.SlackClient.BotUserID
}
