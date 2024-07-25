package prometheus_client

import (
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
	"context"
	"crypto/tls"
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
	prom_client          prom_api.Client
	mimir_client         prom_api.Client
	single_tenant_client prom_api.Client
	http_client          *http.Client
	https_client         *http.Client
	v1api_prom_client    prom_api_v1.API
	v1api_mimir_client   prom_api_v1.API
	v1api_single_tenant  prom_api_v1.API
}

func GetPrometheusClient() *PrometheusClient {
	log.Info().Msgf("client: configuring prometheus client")
	prom_client, err := prom_api.NewClient(prom_api.Config{
		Address: vars.PrometheusUrl,
	})
	if err != nil {
		log.Fatal().Msgf("client: error occured while creating prometheus client: %v", err.Error())
	}
	v1api_prom_client := prom_api_v1.NewAPI(prom_client)
	log.Info().Msgf("client: prometheus client configured, server url [%v]", vars.PrometheusUrl)
	mimir_client, err := prom_api.NewClient(prom_api.Config{
		Address: vars.MimirUrl,
	})
	if err != nil {
		log.Fatal().Msgf("client: error occured while creating mimir client: %v", err.Error())
	}
	v1api_mimir_client := prom_api_v1.NewAPI(mimir_client)
	log.Info().Msgf("client: mimir client configured, server url [%v]", vars.MimirUrl)
	return &PrometheusClient{
		prom_client:        prom_client,
		mimir_client:       mimir_client,
		v1api_prom_client:  v1api_prom_client,
		v1api_mimir_client: v1api_mimir_client,
		http_client:        http.DefaultClient,
		https_client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (prom *PrometheusClient) SetSingleTenantClient(serverUrl string) error {
	single_tenant_client, err := prom_api.NewClient(prom_api.Config{
		Address: serverUrl,
	})
	if err != nil {
		log.Error().Msgf("client: prometheus client could not be configured for server url [%v] %v", serverUrl, err.Error())
		return err
	}
	v1api_single_tenant := prom_api_v1.NewAPI(single_tenant_client)
	prom.single_tenant_client = single_tenant_client
	prom.v1api_single_tenant = v1api_single_tenant
	return nil
}

func (prom *PrometheusClient) ExecuteQuery(query string) (string, error) {
	log.Debug().Msgf("client: executing prometheus query")
	result, warnings, err := prom.v1api_prom_client.Query(context.Background(), query, time.Now())
	if err != nil {
		log.Error().Msgf("client: error occured during exeucitng prometheus query: %v", err.Error())
		return "", err
	}
	if len(warnings) > 0 {
		log.Warn().Msgf("client: warnings appeared during executing prometheus query: %v", warnings)
	}
	return result.String(), nil
}

func (prom *PrometheusClient) ExecuteQueryMimir(query string) (string, error) {
	log.Debug().Msgf("client: executing mimir prometheus query")
	result, warnings, err := prom.v1api_mimir_client.Query(context.Background(), query, time.Now())
	if err != nil {
		log.Error().Msgf("client: error occured during exeucitng mimir query: %v", err.Error())
		return "", err
	}
	if len(warnings) > 0 {
		log.Warn().Msgf("client: warnings appeared during executing mimir query: %v", warnings)
	}
	return result.String(), nil
}

func (prom *PrometheusClient) ExecuteQuerySingleTenant(query string) (string, error) {
	log.Debug().Msgf("client: executing single tenant prometheus query")
	result, warnings, err := prom.v1api_single_tenant.Query(context.Background(), query, time.Now())
	if err != nil {
		log.Error().Msgf("client: error occured during exeucitng single tenant query: %v", err.Error())
		return "", err
	}
	if len(warnings) > 0 {
		log.Warn().Msgf("client: warnings appeared during executing single tenant query: %v", warnings)
	}
	return result.String(), nil
}

func (prom *PrometheusClient) GetDistinctMetricsAndUsage(filter string, singleTenantTarget ...string) (err error) {
	var result string
	query := `count by (__name__)({__name__=~".+"})`
	if len(singleTenantTarget) > 0 {
		for _, singleTenantTarget := range singleTenantTarget {
			log.Debug().Msgf("client: prometheus getting distinct metrics and usage, filer [%v] prometheus single tenant target [%v]", filter, singleTenantTarget)
			serverUrl := fmt.Sprintf(vars.SINGLE_TENANT_PROM_URL_FORMAT, singleTenantTarget)
			err = prom.SetSingleTenantClient(serverUrl)
			if err != nil {
				log.Error().Msgf("clinet: prometheus target [%v] failed to set client for server url [%v] %v", singleTenantTarget, serverUrl, err.Error())
				return err
			}
			result, err = prom.ExecuteQuerySingleTenant(query)
			if err != nil {
				log.Error().Msgf("clinet: prometheus target [%v] failed to execute query %v", singleTenantTarget, err.Error())
				return err
			}
			err = prom.filterAndWriteOutMetrics(result, filter, singleTenantTarget)
			if err != nil {
				return err
			}
		}
	} else {
		log.Debug().Msgf("client: prometheus getting distinct metrics and usage, filer [%v]", filter)
		result, err = prom.ExecuteQuery(query)
		if err != nil {
			log.Error().Msgf("clinet: prometheus failed to execute query %v", err.Error())
			return err
		}
		err = prom.filterAndWriteOutMetrics(result, filter, vars.Environment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (prom *PrometheusClient) filterAndWriteOutMetrics(result string, filter string, fileSuffix string) error {
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
		return err
	}
	fileNameFormat := "output/metric_usage_%s.json"
	fileName := fmt.Sprintf(fileNameFormat, fileSuffix)
	err = os.WriteFile(fileName, metricsJSON, 0644)
	if err != nil {
		log.Error().Msgf("client: error writing metrics to file: %v", err)
		return err
	}
	log.Info().Msgf("client: metrics written to file: %s\n", fileName)
	return nil
}

func (prom *PrometheusClient) GetDistinctMetricsAndUsageMimir(filter string) {
	log.Debug().Msgf("client: prometheus getting distinct metrics and usage, filer [%v]", filter)
	query := `count by (__name__)({__name__=~".+"})`
	result, err := prom.ExecuteQueryMimir(query)
	if err != nil {
		log.Error().Msgf("Failed to execute query %v", err.Error())
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
	fileNameSuffix := "mimir"
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
	log.Debug().Msgf("client: prometheus getting configured alerts from target list %v", prometheusTargets)
	alertingRules := map[string][]entity.AlertingRule{}
	for _, prometheusTarget := range prometheusTargets {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/rules", prometheusTarget), nil)
		if err != nil {
			log.Error().Msgf("client: error occurred while creating HTTP request: %v", err.Error())
			return nil, err
		}
		resp, err := prom.https_client.Do(req)
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

/*
Get all the related alerts of metric list from all existing proemtheus servers,
the returned value is in the format of the following mapping chain:
map[metric name][prometheus url][rule name] -> rule query
(4 dimension matrix)
*/
func (prom *PrometheusClient) GetMetricsAlertsFromAllEnvs(metrics []entity.ExportedMetric) (map[string]map[string]map[string]string, []error) {
	log.Debug().Msgf("client: prometheus getting related alerts from metrics list - all prometheus environments")
	errList := []error{}
	alertMapping := map[string]map[string]map[string]string{}
	for _, envs := range vars.PrometheusServers {
		for _, promeUrl := range envs {
			log.Debug().Msgf("client: prometheus getting related alerts from metrics list - prometheus target [%v]", promeUrl)
			alertingRules, err := prom.GetConfiguredAlerts(promeUrl)
			if err != nil {
				log.Error().Msgf("client: could not perform alert search for prometheus target [%v] %v", promeUrl, err.Error())
				errList = append(errList, err)
				continue
			}
			for _, metric := range metrics {
				if _, hasKey := alertMapping[metric.Name]; !hasKey {
					alertMapping[metric.Name] = map[string]map[string]string{}
				}
				if _, hasKey := alertMapping[metric.Name][promeUrl]; !hasKey {
					alertMapping[metric.Name][promeUrl] = map[string]string{}
				}
				for _, alertRuleList := range alertingRules {
					for _, rule := range alertRuleList {
						if strings.Contains(rule.Query, metric.Name) {
							alertMapping[metric.Name][promeUrl][rule.Name] = rule.Query
						}
					}
				}
			}
		}
	}
	if len(errList) > 0 {
		return alertMapping, errList
	}
	return alertMapping, nil
}

// func (prom *PrometheusClient) GetMetricOwningTeam(metric string) (string, error) {
// 	query := fmt.Sprintf("count by (owner) (%v * on (srv) group_left(owner) (group by (srv, owner) (armis_pg_multitenant_management_service_owner_creation_time)))", metric)
// 	fmt.Println(query)
// 	str, err := prom.ExecuteQuery(query)
// 	if err != nil {
// 		log.Error().Msgf("client: prometheus could not get metric owning team %v", err.Error())
// 		return "", err
// 	}

// 	return str, nil
// }
