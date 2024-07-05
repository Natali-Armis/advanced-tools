package k8s_client

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	log.Info().Msgf("configuring k8s client")
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Info().Msgf("in-cluster config not found, falling back to out-of-cluster config")
		var kubeconfigPath *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfigPath = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		configOverrides := &clientcmd.ConfigOverrides{
			CurrentContext: contextName,
		}
		kubeConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfigPath},
			configOverrides).ClientConfig()
		if err != nil {
			log.Fatal().Msgf("unable to load Kubernetes client config, %v", err)
		}
	} else {
		log.Info().Msgf("in-cluster config loaded successfully")
	}
	k8sClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal().Msgf("unable to create Kubernetes client, %v", err)
	}
	log.Info().Msgf("k8s client configured, context: %v", contextName)
	return &K8sClient{
		k8sClient: k8sClient,
		config:    kubeConfig,
		ctx:       contextName,
	}
}

func (client *K8sClient) GetNodes() (*v1.NodeList, error) {
	nodes, err := client.k8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error().Msgf("unable to list kubernetes nodes, %v", err)
		return nil, err
	}
	return nodes, nil
}

func (client *K8sClient) GetNodeVersionAndLabels(nodeName string) (string, string, error) {
	node, err := client.k8sClient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Error().Msgf("could not get node %v metada: %v", nodeName, err.Error())
		return "", "", err
	}
	version := node.Status.NodeInfo.KubeletVersion
	labels := node.Labels["armis.com/services"]
	return version, labels, nil
}
