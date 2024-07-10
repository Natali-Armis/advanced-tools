package vars

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
	PROMETHEUS_SERVER_TYPE              = "SERVER_TYPE"
	PROMETHEUS_URL                      = "PROMETHEUS_URL"
	ALERT_MANAGER_URL                   = "ALERT_MANAGER_URL"
	ENVIRONMENT                         = "ENVIRONMENT"
	LOG_LEVEL                           = "LOG_LEVEL"
	AWS_PROFILE                         = "AWS_PROFILE"
	AWS_REGION                          = "AWS_REGION"
	SLACK_AUTH_TOKEN                    = "SLACK_AUTH_TOKEN"
	SLACK_UPGRADE_NOTIFICATIONS_CHANNEL = "SLACK_UPGRADE_NOTIFICATIONS_CHANNEL"

	// env vars default
	PROMETHEUS_SERVER_TYPE_DEFAULT = FEDERATION
	ENVIRONMENT_DEFAULT            = DEV
	LOG_LEVEL_DEFAULT              = DEBUG
	AWS_PROFILE_DEFAULT            = PROFILE_DEFAULT
	AWS_REGION_DEFAULT             = US_EAST_1
	ALERT_MANAGER_URL_DEFAULT      = "http://prometheus-federation-prod.armis.internal:9093"

	// misc
	ASG_REQUIRED_LABEL            = "k8s_armis_com"
	INGRESS_NAMESPACE             = "ingress-nginx"
	INGRESS_LABEL_SELECTOR        = "app.kubernetes.io/component"
	INGRESS_LABEL_SELECTOR_VALUE  = "controller"
	INGRESS_NODE_SELECTOR_MATCHER = "armis.com/node-version"
	ALERTMANAGER_SILENCE_ROUT     = "/api/v2/silences"
	MAX_SLACK_MESSAGE_SIZE        = 4000

	// alert manager silencing
	ALERT_NAME_LABEL                   = "alertname"
	ALERT_ENV_LABEL                    = "env"
	PIPELINE_TARGET_DOWN               = "pipeline target down"
	PIPELINE_CONTAINER_RESTARTING      = "pipeline container restarting tenant"
	SERVER_INTEGRATION_CRASHED         = "server integration crashed on tenant"
	CONTAINER_CRASHED_MULTIPLE_TENANTS = "container crashed on multiple tenants"

	// unhealthy pod states
	IMAGE_PULL_ERROR              = "ImagePullBackOff"
	ERR_IMAGE_PULL                = "ErrImagePull"
	CRASH_LOOP_BACKOFF            = "CrashLoopBackOff"
	RUN_CONTAINER_ERROR           = "RunContainerError"
	CREATE_CONTAINER_CONFIG_ERROR = "CreateContainerConfigError"
	CREATE_CONTAINER_ERROR        = "CreateContainerError"

	// upgrade manager stages
	PRE_PROCESS       = "PRE_PROCESS_VALIDATIONS"
	PENDING_APPROVAL  = "PENDING_APPROVAL"
	EDITING_TF_FILES  = "EDITING_TF_FILES"
	APPLYING_TF       = "APPLYING_TF"
	ROLLING_EKS_NODES = "ROLLING_EKS_NODES"
	POST_PROCESS      = "POST_PROCESS_VALIDATIONS"

	// upgrade
	UPGEADE_INITIALIZE_SUBSTR  = `kubernetes-manager upgrade cluster`
	UPGRADE_INITIALIZE_PATTERN = `kubernetes-manager upgrade cluster (\d+\.\d+)`
)
