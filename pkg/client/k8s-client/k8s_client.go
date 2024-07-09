package k8s_client

import (
	"advanced-tools/pkg/vars"
	"context"
	"flag"
	"path/filepath"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClient struct {
	k8sClient *kubernetes.Clientset
	config    *rest.Config
	ctx       string
}

func GetK8sClient(contextName string) *K8sClient {
	log.Info().Msgf("client: configuring k8s client")
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Info().Msgf("client: in-cluster config not found, falling back to out-of-cluster config")
		var kubeconfigPath *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfigPath = flag.String("client: kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfigPath = flag.String("client: kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		configOverrides := &clientcmd.ConfigOverrides{
			CurrentContext: contextName,
		}
		kubeConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfigPath},
			configOverrides).ClientConfig()
		if err != nil {
			log.Fatal().Msgf("client: unable to load Kubernetes client config, %v", err)
		}
	} else {
		log.Info().Msgf("client: in-cluster config loaded successfully")
	}
	k8sClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal().Msgf("client: unable to create Kubernetes client, %v", err)
	}
	log.Info().Msgf("client: k8s client configured, context: %v", contextName)
	return &K8sClient{
		k8sClient: k8sClient,
		config:    kubeConfig,
		ctx:       contextName,
	}
}

func (client *K8sClient) GetNodes() (*core_v1.NodeList, error) {
	nodes, err := client.k8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error().Msgf("client: unable to list kubernetes nodes, %v", err)
		return nil, err
	}
	return nodes, nil
}

func (client *K8sClient) GetNodeVersionAndLabel(nodeName string) (string, string, error) {
	node, err := client.k8sClient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Error().Msgf("client: could not get node %v metada: %v", nodeName, err.Error())
		return "", "", err
	}
	version := node.Status.NodeInfo.KubeletVersion
	label := node.Labels["armis.com/services"]
	return version, label, nil
}

func (client *K8sClient) GetDeploymentsByLabelSlector(labelSelectorMap map[string]string) ([]v1.Deployment, error) {
	var selector labels.Selector
	for key, value := range labelSelectorMap {
		selector = labels.SelectorFromSet(labels.Set{key: value})
	}

	deployments, err := client.k8sClient.AppsV1().Deployments(vars.INGRESS_NAMESPACE).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		log.Error().Msgf("client: unable to list deployments with label selector %v, %v", labelSelectorMap, err)
		return nil, err
	}
	return deployments.Items, nil
}

func (client *K8sClient) EditIngressDeploymentToMatchVersionLabel(deployment *v1.Deployment, newVersion string) error {
	if deployment.Spec.Template.Spec.NodeSelector == nil {
		deployment.Spec.Template.Spec.NodeSelector = make(map[string]string)
	}
	deployment.Spec.Template.Spec.NodeSelector[vars.INGRESS_NODE_SELECTOR_MATCHER] = newVersion
	_, err := client.k8sClient.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		log.Error().Msgf("client: unable to update deployment %v, %v", deployment.Name, err)
		return err
	}
	return nil
}

func (client *K8sClient) VerifyIngressDeploymentToMatchVersionLabel(deployment *v1.Deployment, newVersion string) (bool, error) {
	if deployment.Spec.Template.Spec.NodeSelector == nil {
		return false, nil
	}
	currentVersion, exists := deployment.Spec.Template.Spec.NodeSelector[vars.INGRESS_NODE_SELECTOR_MATCHER]
	if !exists {
		return false, nil
	}
	log.Debug().Msgf("client: ingress deployment [%v] has version label [%v]", deployment.Name, currentVersion)
	if currentVersion == newVersion {
		return true, nil
	}
	return false, nil
}
