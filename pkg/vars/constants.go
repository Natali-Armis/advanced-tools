package vars

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

const (
	// server types
	FEDERATION = "federation"
	BACKEND    = "backend"
	K8_PROXY   = "k8s-proxy"

	// environments
	DEV             = "dev1"
	QA              = "qa1"
	OPERATIONS_TEST = "operations-test"
	DEMO            = "demo1"
	PROD            = "prod"
	PROD4           = "prod4"
	PROD5           = "prod5"
	PROD7           = "prod7"
	PROD8           = "prod8"
	PROD9           = "prod9"
	OPERATIONS      = "operations"

	// aws profiles
	PROFILE_DEFAULT = "default"
	PROFILE_GOV     = "gov"

	// aws regions
	US_EAST_1      = "us-east-1"      // dev1, qa1, demo1, operations_test, operations, prod
	EU_CENTRAL_1   = "eu-central-1"   // prod4
	AP_SOUTHEAST_1 = "ap-southeast-1" // prod7
	AP_SOUTHEAST_2 = "ap-southeast-2" // prod5
	US_GOV_EAST_1  = "us-gov-east-1"  // prod8
	AP_SOUTH_1     = "ap-south-1"     // prod9

	// port
	PROMETHEUS_PORT = "9090"

	// log levels
	ERROR = "erro"
	WARN  = "warn"
	INFO  = "info"
	DEBUG = "debug"

	// env vars names
	PROMETHEUS_SERVER_TYPE = "SERVER_TYPE"
	PROMETHEUS_URL         = "PROMETHEUS_URL"
	ENVIRONMENT            = "ENVIRONMENT"
	LOG_LEVEL              = "LOG_LEVEL"
	AWS_PROFILE            = "AWS_PROFILE"
	AWS_REGION             = "AWS_REGION"
	SLACK_AUTH_TOKEN       = "SLACK_AUTH_TOKEN"

	// env vars default
	PROMETHEUS_SERVER_TYPE_DEFAULT = FEDERATION
	ENVIRONMENT_DEFAULT            = DEV
	LOG_LEVEL_DEFAULT              = DEBUG
	AWS_PROFILE_DEFAULT            = PROFILE_DEFAULT
	AWS_REGION_DEFAULT             = US_EAST_1

	// misc
	ASG_REQUIRED_LABEL            = "k8s_armis_com"
	INGRESS_NAMESPACE             = "ingress-nginx"
	INGRESS_LABEL_SELECTOR        = "app.kubernetes.io/component"
	INGRESS_LABEL_SELECTOR_VALUE  = "controller"
	INGRESS_NODE_SELECTOR_MATCHER = "armis.com/node-version"
)

var (
	PrometheusServerType = ""
	Environment          = ""
	PrometheusUrl        = ""
	LogLevel             = ""
	AwsProfile           = ""
	AwsRegion            = ""

	SlackAuthToken = ""

	ClusterToRegionMapper map[string]string
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
	log.Debug().Msgf("config: %v [%v]", PROMETHEUS_URL, PrometheusUrl)
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
	SlackAuthToken = os.Getenv(SLACK_AUTH_TOKEN)
	if len(SlackAuthToken) == 0 {
		log.Fatal().Msgf("environment variable %v must be defined", SLACK_AUTH_TOKEN)
	}
	log.Debug().Msgf("config: %v was set to non empty value", SLACK_AUTH_TOKEN)
}
