package client

import (
	alertmanager "advanced-tools/pkg/client/alertmanager-client"
	aws "advanced-tools/pkg/client/aws-client"
	grafana "advanced-tools/pkg/client/grafana-client"
	k8s "advanced-tools/pkg/client/k8s-client"
	prom "advanced-tools/pkg/client/prometheus-client"
	slack "advanced-tools/pkg/client/slack-client"
	"advanced-tools/pkg/vars"
)

type Client struct {
	PrometheusClient   *prom.PrometheusClient
	AlertManagerClient *alertmanager.AlertManagerClient
	K8sClient          *k8s.K8sClient
	AwsClient          *aws.AwsClient
	SlackClient        *slack.SlackClient
	GrafanaClient      *grafana.GrafanaClient
}

func GetClient() *Client {
	return &Client{
		PrometheusClient:   prom.GetPrometheusClient(),
		AlertManagerClient: alertmanager.GetAlertManagerClient(),
		K8sClient:          k8s.GetK8sClient(vars.Environment),
		AwsClient:          aws.GetAwsClient(vars.AwsProfile, vars.AwsRegion),
		SlackClient:        slack.GetSlackClient(vars.SlackAuthToken),
		GrafanaClient:      grafana.GetGrafanaClient(vars.GrafanaToken, vars.GRAFANA_URL),
	}
}
