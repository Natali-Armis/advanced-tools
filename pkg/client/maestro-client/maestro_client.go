package maestro_client

import (
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type MaestroClient struct {
	https_client *http.Client
}

func GetMaestroClient() *MaestroClient {
	return &MaestroClient{
		https_client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (maestro *MaestroClient) GetTasks(params map[string]string) (*entity.MaestroTaskResponse, error) {
	url := fmt.Sprintf("%vtasks?", vars.MaestroUrl)
	for k, v := range params {
		if k == "length" {
			continue
		}
		url = fmt.Sprintf("%v%v=%v&", url, k, v)
	}
	if len, hasKey := params["length"]; hasKey {
		url = fmt.Sprintf("%vlength=%v", url, len)
	} else {
		url = fmt.Sprintf("%vlength=%v", url, vars.MaestroMaxResultsLimit)
	}
	resp, err := maestro.https_client.Get(url)
	if err != nil {
		log.Error().Msgf("client: maestro could not perform get tasks request %v", err.Error())
		return nil, err
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("client: maestro client could not pars response for get tasks %v", err.Error())
		return nil, err
	}
	maestroTaskResponse := &entity.MaestroTaskResponse{}
	err = json.Unmarshal(bytes, maestroTaskResponse)
	if err != nil {
		log.Error().Msgf("client: maestro client could not unmarshall response for get tasks %v", err.Error())
		return nil, err
	}
	return maestroTaskResponse, nil
}

func (maestro *MaestroClient) GetTenants(params map[string]string) (*entity.MaestroTenantResponse, error) {
	url := fmt.Sprintf("%vtenants?", vars.MaestroUrl)
	for k, v := range params {
		if k == "length" {
			continue
		}
		url = fmt.Sprintf("%v%v=%v&", url, k, v)
	}
	if len, hasKey := params["length"]; hasKey {
		url = fmt.Sprintf("%vlength=%v", url, len)
	} else {
		url = fmt.Sprintf("%vlength=%v", url, vars.MaestroMaxResultsLimit)
	}
	// fmt.Println(url)
	resp, err := maestro.https_client.Get(url)
	if err != nil {
		log.Error().Msgf("client: maestro could not perform get tenants request %v", err.Error())
		return nil, err
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("client: maestro client could not pars response for get tenants %v", err.Error())
		return nil, err
	}
	// fmt.Println(string(bytes))
	maestroTenantResponse := &entity.MaestroTenantResponse{}
	err = json.Unmarshal(bytes, maestroTenantResponse)
	if err != nil {
		log.Error().Msgf("client: maestro client could not unmarshall response for get tenants %v", err.Error())
		return nil, err
	}
	return maestroTenantResponse, nil
}
