package grafana_client

import (
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type GrafanaClient struct {
	apiKey    string
	serverUrl string
	client    *http.Client
}

func GetGrafanaClient(apiKey string, serverUrl string) *GrafanaClient {
	return &GrafanaClient{
		apiKey:    apiKey,
		serverUrl: serverUrl,
		client:    &http.Client{},
	}
}

func (client *GrafanaClient) GetAllDashboards() ([]entity.DashboardSearchResult, error) {
	url := fmt.Sprintf("%v%v", client.serverUrl, vars.GRAFANA_ALL_DASHBAORDS_ROUT)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Msgf("client: could not create request for grafana server %v", err.Error())
		return nil, err
	}
	addHeaders(req, map[string]string{
		"Authorization": "Bearer " + client.apiKey,
	})
	resp, err := client.client.Do(req)
	if err != nil {
		log.Error().Msgf("client: could not perform request to grafana server %v", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("client: could not perform request to grafana server %v", err.Error())
		return nil, err
	}
	var results []entity.DashboardSearchResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Error().Msgf("client: could not unmarshal results from grafana server %v", err.Error())
		return nil, err
	}
	return results, nil
}

func (client *GrafanaClient) GetDashboardByUID(uid string) (entity.Dashboard, error) {
	url := fmt.Sprintf("%v%v%v", client.serverUrl, vars.GRAFANA_SINGLE_DASHBOARD_ROUT, uid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Msgf("client: could not create request for grafana server %v", err.Error())
		return entity.Dashboard{}, err
	}
	addHeaders(req, map[string]string{
		"Authorization": "Bearer " + client.apiKey,
	})
	resp, err := client.client.Do(req)
	if err != nil {
		log.Error().Msgf("client: could not perform request to grafana server %v", err.Error())
		return entity.Dashboard{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("client: could not perform request to grafana server %v", err.Error())
		return entity.Dashboard{}, err
	}
	var dashboardResponse entity.DashboardResponse
	err = json.Unmarshal(body, &dashboardResponse)
	if err != nil {
		log.Error().Msgf("client: could not unmarshal dashboard response from grafana server %v", err.Error())
		return entity.Dashboard{}, err
	}
	return dashboardResponse.Dashboard, nil
}

func (client *GrafanaClient) FindMetricInDashboards(metricName string) ([]string, error) {
	dashboards, err := client.GetAllDashboards()
	if err != nil {
		return nil, err
	}
	containingDashabords := []string{}
	for _, dashboard := range dashboards {
		dash, err := client.GetDashboardByUID(dashboard.UID)
		if err != nil {
			log.Error().Msgf("client: could not get dashboard by UID %v", err.Error())
			continue
		}
		if containsMetric(dash, metricName) {
			containingDashabords = append(containingDashabords, dash.Title)
		}
	}
	return containingDashabords, nil
}

func (client *GrafanaClient) FindMetricsInDashboards(metrics []entity.ExportedMetric) (map[string][]string, error) {
	dashboards, err := client.GetAllDashboards()
	if err != nil {
		return nil, err
	}
	containingDashabords := map[string][]string{}
	for _, dashboard := range dashboards {
		dash, err := client.GetDashboardByUID(dashboard.UID)
		if err != nil {
			log.Error().Msgf("client: could not get dashboard by UID %v", err.Error())
			continue
		}
		for _, metric := range metrics {
			metricName := metric.Name
			if _, hasKey := containingDashabords[metricName]; !hasKey {
				containingDashabords[metricName] = []string{}
			}
			if containsMetric(dash, metricName) {
				containingDashabords[metricName] = append(containingDashabords[metricName], dash.Title)
			}
		}
	}
	return containingDashabords, nil
}

func containsMetric(dashboard entity.Dashboard, metricName string) bool {
	dashboardJSON, err := json.Marshal(dashboard)
	if err != nil {
		log.Error().Msgf("client: could not marshal dashboard %v", err.Error())
		return false
	}
	return strings.Contains(string(dashboardJSON), metricName)
}

func addHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
