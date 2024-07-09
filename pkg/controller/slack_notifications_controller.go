package controller

import (
	"advanced-tools/pkg/client"
)

type SlackNotificationController struct {
	clients *client.Client
}

func GetSlackNotificationController(clients *client.Client) *SlackNotificationController {
	return &SlackNotificationController{
		clients: clients,
	}
}
