package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"

	"encoding/json"
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

	// clients.PrometheusClient.GetDistinctMetricsAndUsageMimir("armis")

	// err := clients.PrometheusClient.GetDistinctMetricsAndUsage("armis", vars.SingleTenantTargets...)
	// if err != nil {
	// 	return
	// }

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
	// outputFileMetricsDashabords := "output/metrics_dashboards.json"
	// err = os.WriteFile(outputFileMetricsDashabords, bytes, 0777)
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
	// outputFileMetricsAlerts := "output/metrics_alerts.json"
	// err = os.WriteFile(outputFileMetricsAlerts, bytes, 0777)
	// if err != nil {
	// 	log.Error().Msgf("could not perform metrics alerts writing to output file [%v] %v", outputFile, err.Error())
	// 	return
	// }

	// outputMetricsWithoutAlertAndDashboard(outputFileMetricsDashabords, outputFileMetricsAlerts)

	// for _, metric := range exportedMetrics {
	// 	clients.PrometheusClient.GetMetricOwningTeam(metric.Name)
	// }

}

func outputMetricsWithoutAlertAndDashboard(metricsDashboardsFile string, metricsAlertsFile string) error {
	metricsWithoutAlertsAndDashboards := []string{}
	bytes, err := os.ReadFile(metricsAlertsFile)
	if err != nil {
		log.Error().Msgf("could not open file [%v] %v", metricsAlertsFile, err.Error())
		return err
	}
	var metricsAlerts map[string]map[string]map[string]string
	err = json.Unmarshal(bytes, &metricsAlerts)
	if err != nil {
		log.Error().Msgf("could not unmarshall input from %v %v", metricsAlertsFile, err.Error())
		return err
	}

	bytes, err = os.ReadFile(metricsDashboardsFile)
	if err != nil {
		log.Error().Msgf("could not open file [%v] %v", metricsAlertsFile, err.Error())
		return err
	}
	var metricsDashboards map[string][]string
	err = json.Unmarshal(bytes, &metricsDashboards)
	if err != nil {
		log.Error().Msgf("could not unmarshall input from %v %v", metricsDashboardsFile, err.Error())
		return err
	}

	for metric, dashabords := range metricsDashboards {
		if len(dashabords) > 0 {
			continue
		}
		promUrls := metricsAlerts[metric]
		if len(promUrls) > 0 {
			continue
		}
		metricsWithoutAlertsAndDashboards = append(metricsWithoutAlertsAndDashboards, metric)
	}

	bytes, err = json.MarshalIndent(metricsWithoutAlertsAndDashboards, "", "    ")
	if err != nil {
		log.Error().Msgf("could not marshall input %v", err.Error())
		return err
	}

	outputFile := "output/metrics_without_alerts_and_dashabords.json"
	err = os.WriteFile(outputFile, bytes, 0777)
	if err != nil {
		log.Error().Msgf("could not perform metrics writing to output file [%v] %v", outputFile, err.Error())
		return err
	}
	return nil
}
