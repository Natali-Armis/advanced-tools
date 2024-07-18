package vars

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

var (
	PrometheusServerType string
	Environment          string
	PrometheusUrl        string
	AlertManagerUrl      string
	LogLevel             string
	AwsProfile           string
	AwsRegion            string

	GrafanaToken                     string
	SlackAuthToken                   string
	SlackUpgradeNotificationsChannel string

	ClusterToRegionMapper        map[string]string
	AlertsToSilenceDuringUpgrade []string
	AlertsTeamsToNotSilence      []string
	AwsInstancesCodes            map[int32]string

	SingleTenantTargets []string
)

func init() {
	LogLevel = os.Getenv(LOG_LEVEL)
	if len(LogLevel) == 0 {
		LogLevel = LOG_LEVEL_DEFAULT
	}
	log.Debug().Msgf("config: %v [%v]", LOG_LEVEL, LogLevel)
	PrometheusServerType = os.Getenv(PROMETHEUS_SERVER_TYPE)
	if len(PrometheusServerType) == 0 {
		PrometheusServerType = PROMETHEUS_SERVER_TYPE_DEFAULT
	}
	log.Debug().Msgf("config: %v [%v]", PROMETHEUS_SERVER_TYPE, PrometheusServerType)
	Environment = os.Getenv(ENVIRONMENT)
	if len(Environment) == 0 {
		Environment = ENVIRONMENT_DEFAULT
	}
	log.Debug().Msgf("config: %v [%v]", ENVIRONMENT, Environment)
	PrometheusUrl = os.Getenv(PROMETHEUS_URL)
	if len(PrometheusUrl) == 0 {
		PrometheusUrl = fmt.Sprintf("http://prometheus-%v-%v.armis.internal:%v", PrometheusServerType, Environment, PROMETHEUS_PORT)
	}
	PrometheusServerFormats := map[string]string{
		FEDERATION: "http://prometheus-%v-%v.armis.internal:%v",
		BACKEND:    "http://prometheus-%v-%v.armis.internal:%v",
		K8_PROXY:   "http://prometheus-%v.%v.k8s.armis.com:%v",
	}
	PrometheusServerEnvs := map[string][]string{
		FEDERATION: []string{DEV, QA, DEMO, PROD, PROD4, PROD5, PROD7, PROD8, PROD9},
		BACKEND:    []string{DEV, QA, DEMO, PROD, PROD4, PROD5, PROD7, PROD8, PROD9},
		K8_PROXY:   []string{OPERATIONS, OPERATIONS_TEST, DEV, QA, DEMO, PROD, PROD4, PROD5, PROD7, PROD8, PROD9},
	}

	log.Debug().Msgf("config: %v [%v]", PROMETHEUS_URL, PrometheusUrl)
	AlertManagerUrl = os.Getenv(ALERT_MANAGER_URL)
	if len(AlertManagerUrl) == 0 {
		AlertManagerUrl = ALERT_MANAGER_URL_DEFAULT
	}
	log.Debug().Msgf("config: %v [%v]", ALERT_MANAGER_URL, AlertManagerUrl)
	AwsProfile = os.Getenv(AWS_PROFILE)
	if len(AwsProfile) == 0 {
		AwsProfile = AWS_PROFILE_DEFAULT
	}
	log.Debug().Msgf("config: %v [%v]", AWS_PROFILE, AwsProfile)
	AwsRegion = os.Getenv(AWS_REGION)
	if len(AwsRegion) == 0 {
		AwsRegion = AWS_REGION_DEFAULT
	}
	log.Debug().Msgf("config: %v [%v]", AWS_REGION, AwsRegion)
	ClusterToRegionMapper = map[string]string{
		DEV:             US_EAST_1,
		QA:              US_EAST_1,
		DEMO:            US_EAST_1,
		OPERATIONS_TEST: US_EAST_1,
		OPERATIONS:      US_EAST_1,
		PROD:            US_EAST_1,
		PROD4:           EU_CENTRAL_1,
		PROD5:           AP_SOUTHEAST_2,
		PROD7:           AP_SOUTHEAST_1,
		PROD8:           US_GOV_EAST_1,
		PROD9:           AP_SOUTH_1,
	}
	AlertsToSilenceDuringUpgrade = []string{
		PIPELINE_TARGET_DOWN,
		PIPELINE_CONTAINER_RESTARTING,
		SERVER_INTEGRATION_CRASHED,
		SERVER_INTEGRATION_CRASHED_MULTIPLE,
		CONTAINER_CRASHED_MULTIPLE_TENANTS,
		CONTAINER_GOT_OOM,
	}
	AlertsTeamsToNotSilence = []string{
		DEVOPS_INFRA,
		SOLUTIONS,
		APPS_INFRA,
		DATA_INFRA,
		STREAMING_INFRA,
	}
	AwsInstancesCodes = map[int32]string{
		0:  "pending",
		16: "running",
		32: "shutting-down",
		64: "stopping",
	}
	GrafanaToken = os.Getenv(GRAFANA_TOKEN)
	if len(GrafanaToken) == 0 {
		log.Fatal().Msgf("environment variable %v must be defined", GRAFANA_TOKEN)
	}
	SlackAuthToken = os.Getenv(SLACK_AUTH_TOKEN)
	if len(SlackAuthToken) == 0 {
		log.Fatal().Msgf("environment variable %v must be defined", SLACK_AUTH_TOKEN)
	}
	log.Debug().Msgf("config: %v was set to non empty value", SLACK_AUTH_TOKEN)
	SlackUpgradeNotificationsChannel = os.Getenv(SLACK_UPGRADE_NOTIFICATIONS_CHANNEL)
	if len(SlackUpgradeNotificationsChannel) == 0 {
		log.Fatal().Msgf("environment variable %v must be defined", SLACK_UPGRADE_NOTIFICATIONS_CHANNEL)
	}
	log.Debug().Msgf("config: %v [%v]", SLACK_UPGRADE_NOTIFICATIONS_CHANNEL, SlackUpgradeNotificationsChannel)
	SingleTenantTargets = []string{
		"nestle",
		"lowes",
		"kaiser",
		"home-depot",
		"nvidia",
		"intel-fab",
		"jaguar-landrover",
		"walmart",
		"tesla",
		"helix",
	}
}
