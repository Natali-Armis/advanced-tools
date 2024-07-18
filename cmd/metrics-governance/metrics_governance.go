package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"
	"advanced-tools/pkg/entity"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
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
	// 	clients.PrometheusClient.GetDistinctMetricsAndUsage("armis", tenant)
	// }

	// dashboards, err := clients.GrafanaClient.GetAllDashboards()
	// if err != nil {
	// 	return
	// }
	// for _, dashabord := range dashboards {
	// 	fmt.Println(dashabord.Title)
	// }

	outputFileName := "output/metric_usage_kaiser.json"
	bytes, err := os.ReadFile(outputFileName)
	if err != nil {
		log.Error().Msgf("could not open file [%v] %v", outputFileName, err.Error())
		return
	}
	var exportedMetrics []entity.ExportedMetric
	err = json.Unmarshal(bytes, &exportedMetrics)
	if err != nil {
		log.Error().Msgf("could not open file [%v] %v", outputFileName, err.Error())
		return
	}

	// metricsDashabords, err := clients.GrafanaClient.FindMetricsInDashboards(exportedMetrics)
	// if err != nil {
	// 	log.Error().Msgf("could not perform metrics search in dashboards %v", err.Error())
	// 	return
	// }

	metricsAlerts, err := clients.PrometheusClient.GetMetricsAlerts(exportedMetrics)
	if err != nil {
		log.Error().Msgf("could not perform metrics search in alerts %v", err.Error())
		return
	}
	for metric, alertList := range metricsAlerts {
		if len(alertList) > 0 {
			fmt.Printf("%v | %v\n", metric, alertList)
		}
	}

}
