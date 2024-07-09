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
