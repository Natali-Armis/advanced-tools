package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"
	"advanced-tools/pkg/manager"
)

var (
	clients *client.Client
)

func init() {
	config.Configure()
	clients = client.GetClient()
}

func main() {
	upgradeManager := manager.GetKubernetesUpgradeManager(clients)
	upgradeManager.Run()
}
