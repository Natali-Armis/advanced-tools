package prometheus_client

import (
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func GetPrometheusClient() *PrometheusClient {
	log.Info().Msgf("client: configuring prometheus client")
	client, err := prom_api.NewClient(prom_api.Config{
		Address: vars.PrometheusUrl,
	})
	if err != nil {
		log.Fatal().Msgf("client: error occured while creating prometheus client: %v", err.Error())
	}
	v1api := prom_api_v1.NewAPI(client)
	log.Info().Msgf("client: prometheus client configured, server url [%v]", vars.PrometheusUrl)
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

func (prom *PrometheusClient) GetDistinctMetricsAndUsage(filter string, singleTenantTarget ...string) {
	query := `count by (__name__)({__name__=~".+"})`
	result, err := prom.ExecuteQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return
	}
	metricsLines := strings.Split(result, "\n")
	var metrics []entity.MetricUsageCount
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
			if strings.Contains(name, filter) {
				metrics = append(metrics, entity.MetricUsageCount{
					Name:  name,
					Value: value,
				})
			}
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
	fileNameFormat := "output/metric_usage_%s.json"
	fileNameSuffix := vars.Environment
	if len(singleTenantTarget) > 0 {
		fileNameSuffix = singleTenantTarget[0]
	}
	fileName := fmt.Sprintf(fileNameFormat, fileNameSuffix)
	err = os.WriteFile(fileName, metricsJSON, 0644)
	if err != nil {
		log.Error().Msgf("client: error writing metrics to file: %v", err)
		return
	}
	log.Info().Msgf("client: metrics written to file: %s\n", fileName)
}

func (prom *PrometheusClient) GetConfiguredAlerts(prometheusTargets ...string) (map[string][]entity.AlertingRule, error) {
	if len(prometheusTargets) == 0 {
		prometheusTargets = []string{vars.PrometheusUrl}
	}
	alertingRules := map[string][]entity.AlertingRule{}
	for _, prometheusTarget := range prometheusTargets {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/rules", prometheusTarget), nil)
		if err != nil {
			log.Error().Msgf("client: error occurred while creating HTTP request: %v", err.Error())
			return nil, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Msgf("client: error occurred while making HTTP request: %v", err.Error())
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			msg := fmt.Sprintf("client: received non-OK HTTP status: %v", resp.StatusCode)
			log.Error().Msgf(msg)
			return nil, fmt.Errorf(msg)
		}
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Msgf("client: error occurred while reading response: %v", err.Error())
			return nil, err
		}
		var result struct {
			Status string `json:"status"`
			Data   struct {
				Groups []struct {
					Name  string `json:"name"`
					File  string `json:"file"`
					Rules []struct {
						State       string            `json:"state"`
						Name        string            `json:"name"`
						Query       string            `json:"query"`
						Annotations map[string]string `json:"annotations"`
						Labels      map[string]string `json:"labels"`
					} `json:"rules"`
				} `json:"groups"`
			} `json:"data"`
		}
		err = json.Unmarshal(bytes, &result)
		if err != nil {
			log.Error().Msgf("client: error occurred while converting response to json: %v", err.Error())
			return nil, err
		}
		for _, group := range result.Data.Groups {
			for _, rule := range group.Rules {
				if _, exists := alertingRules[rule.Name]; !exists {
					alertingRules[rule.Name] = []entity.AlertingRule{}
				}
				alertingRule := entity.AlertingRule{
					Name:        rule.Name,
					Query:       rule.Query,
					Description: rule.Annotations["description"],
					Severity:    rule.Labels["severity"],
				}
				alertingRules[rule.Name] = append(alertingRules[rule.Name], alertingRule)
			}
		}
	}
	return alertingRules, nil
}

func (prom *PrometheusClient) GetMetricsAlerts(metrics []entity.ExportedMetric) (map[string]map[string]string, error) {
	alertingRules, err := prom.GetConfiguredAlerts()
	if err != nil {
		log.Error().Msgf("client: could not perform alert search %v", err.Error())
		return nil, err
	}
	containingAlerts := map[string]map[string]string{}
	for _, metric := range metrics {
		if _, hasKey := containingAlerts[metric.Name]; !hasKey {
			containingAlerts[metric.Name] = map[string]string{}
		}
		for _, alertRuleList := range alertingRules {
			for _, rule := range alertRuleList {
				if strings.Contains(rule.Query, metric.Name) {
					containingAlerts[metric.Name][rule.Name] = rule.Query
				}
			}
		}
	}
	return containingAlerts, nil
}
