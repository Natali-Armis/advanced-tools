package prometheus_client

import (
	"advanced-tools/pkg/vars"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	prom_api "github.com/prometheus/client_golang/api"
	prom_api_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/rs/zerolog/log"
)

type PrometheusClient struct {
	client prom_api.Client
	v1api  prom_api_v1.API
}

type MetricUsageCount struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func GetPrometheusClient() *PrometheusClient {
	log.Info().Msgf("client: configuring prometheus client")
	client, err := prom_api.NewClient(prom_api.Config{
		Address: vars.PrometheusUrl,
	})
	if err != nil {
		log.Fatal().Msgf("client: error occured while creating prometheus client: %v", err.Error())
	}
	v1api := prom_api_v1.NewAPI(client)
	log.Info().Msgf("client: prometheus client configured, server url %v", vars.PrometheusUrl)
	return &PrometheusClient{
		client: client,
		v1api:  v1api,
	}
}

func (prom *PrometheusClient) ExecuteQuery(query string) (string, error) {
	result, warnings, err := prom.v1api.Query(context.Background(), query, time.Now())
	if err != nil {
		log.Error().Msgf("client: error occured during exeucitng prometheus query: %v", err.Error())
		return "", err
	}
	if len(warnings) > 0 {
		log.Warn().Msgf("client: warnings appeared during executing prometheus query: %v", warnings)
	}
	return result.String(), nil
}

func (prom *PrometheusClient) GetDistinctMetricsAndUsage() {
	query := `count by (__name__)({__name__=~".+"})`
	result, err := prom.ExecuteQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return
	}
	metricsLines := strings.Split(result, "\n")
	var metrics []MetricUsageCount
	reMetric := regexp.MustCompile(`(.+) => (\d+) @\[.+\]`)
	for _, line := range metricsLines {
		match := reMetric.FindStringSubmatch(line)
		if len(match) == 3 {
			name := strings.TrimSpace(match[1])
			value, err := strconv.Atoi(match[2])
			if err != nil {
				log.Error().Msgf("client: error converting value to integer: %v\n", err)
				continue
			}
			metrics = append(metrics, MetricUsageCount{
				Name:  name,
				Value: value,
			})
		}
	}
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Value > metrics[j].Value
	})
	metricsJSON, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		log.Error().Msgf("client: error marshaling metrics to JSON: %v", err)
		return
	}
	fileName := fmt.Sprintf("metric_usage_%s.json", vars.Environment)
	err = os.WriteFile(fileName, metricsJSON, 0644)
	if err != nil {
		log.Error().Msgf("client: error writing metrics to file: %v", err)
		return
	}
	log.Info().Msgf("client: metrics written to file: %s\n", fileName)
}
