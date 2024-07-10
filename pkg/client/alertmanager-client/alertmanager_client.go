package alertmanager_client

import (
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type AlertManagerClient struct {
	url    string
	client *http.Client
}

func GetAlertManagerClient() *AlertManagerClient {
	return &AlertManagerClient{
		url:    vars.AlertManagerUrl,
		client: &http.Client{},
	}
}

func (client *AlertManagerClient) CreateSilence(silence *entity.Silence) error {
	silenceUrl := fmt.Sprintf("%v%v", vars.AlertManagerUrl, vars.ALERTMANAGER_SILENCE_ROUT)
	silencePayload, err := json.Marshal(silence)
	if err != nil {
		log.Error().Msgf("client: could not create silence %v: %v", silence, err.Error())
		return err
	}
	req, err := http.NewRequest("POST", silenceUrl, bytes.NewBuffer(silencePayload))
	if err != nil {
		log.Error().Msgf("client: could not create http request for silence %v: %v", silence, err.Error())
		return err
	}
	addHeaders(map[string]string{}, req)
	resp, err := client.client.Do(req)
	if err != nil {
		log.Error().Msgf("client: could not create http response for silence %v: %v", silence, err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("failed to create silence, status code: %d", resp.StatusCode)
		log.Error().Msgf("client: %v", err.Error())
		return err
	}
	return nil
}

func addHeaders(headers map[string]string, req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
