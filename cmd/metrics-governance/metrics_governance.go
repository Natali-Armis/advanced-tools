package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"
	"fmt"
)

var (
	clients *client.Client
)

func init() {
	config.Configure()
	clients = client.GetClient()
}

func main() {
	// for _, tenant := range vars.SingleTenantTargets {
	// 	clients.PrometheusClient.GetDistinctMetricsAndUsage(tenant)
	// }
	dashboards, err := clients.GrafanaClient.GetAllDashbaords()
	if err != nil {
		return
	}
	for _, dashabord := range dashboards {
		fmt.Println(dashabord.Title)
	}
}
