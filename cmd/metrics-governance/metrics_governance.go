package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"
)

var (
	clients *client.Client
)

func init() {
	config.Configure()
	clients = client.GetClient()
}

func main() {
	clients.PrometheusClient.GetDistinctMetricsAndUsage()
}
