package grafana_client

import (
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func (client *GrafanaClient) GetAllDashbaords() ([]entity.DashboardSearchResult, error) {
	url := fmt.Sprintf("%v%v", client.serverUrl, vars.GRAFANA_DASHBAORDS_ROUT)
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

func addHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
