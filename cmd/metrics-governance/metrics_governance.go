package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"
	"advanced-tools/pkg/vars"
)

var (
	clients *client.Client
)

func init() {
	config.Configure()
	clients = client.GetClient()
}

func main() {

	err := clients.PrometheusClient.GetDistinctMetricsAndUsage("armis", vars.SingleTenantTargets...)
	if err != nil {
		return
	}

	// outputFileName := "output/metric_usage_kaiser.json"
	// bytes, err := os.ReadFile(outputFileName)
	// if err != nil {
	// 	log.Error().Msgf("could not open file [%v] %v", outputFileName, err.Error())
	// 	return
	// }
	// var exportedMetrics []entity.ExportedMetric
	// err = json.Unmarshal(bytes, &exportedMetrics)
	// if err != nil {
	// 	log.Error().Msgf("could not open file [%v] %v", outputFileName, err.Error())
	// 	return
	// }

	// metricsDashabords, err := clients.GrafanaClient.FindMetricsInDashboards(exportedMetrics)
	// if err != nil {
	// 	log.Error().Msgf("could not perform metrics search in dashboards %v", err.Error())
	// 	return
	// }
	// bytes, err = json.MarshalIndent(metricsDashabords, "", "    ")
	// if err != nil {
	// 	log.Error().Msgf("could not perform metrics dashbords marshaling %v", err.Error())
	// 	return
	// }
	// outputFile := "output/metrics_dashboards.json"
	// err = os.WriteFile(outputFile, bytes, 0777)
	// if err != nil {
	// 	log.Error().Msgf("could not perform metrics alerts writing to output file [%v] %v", outputFile, err.Error())
	// 	return
	// }

	// metricsAlertsFromAllEnvs, errs := clients.PrometheusClient.GetMetricsAlertsFromAllEnvs(exportedMetrics)
	// if errs != nil {
	// 	log.Error().Msgf("could not perform metrics search in alerts")
	// 	for _, err := range errs {
	// 		log.Error().Msg(err.Error())
	// 	}
	// }
	// bytes, err = json.MarshalIndent(metricsAlertsFromAllEnvs, "", "    ")
	// if err != nil {
	// 	log.Error().Msgf("could not perform metrics alerts marshaling %v", err.Error())
	// 	return
	// }
	// outputFile = "output/metrics_alerts.json"
	// err = os.WriteFile(outputFile, bytes, 0777)
	// if err != nil {
	// 	log.Error().Msgf("could not perform metrics alerts writing to output file [%v] %v", outputFile, err.Error())
	// 	return
	// }

	// for _, metric := range exportedMetrics {
	// 	clients.PrometheusClient.GetMetricOwningTeam(metric.Name)
	// }

}
