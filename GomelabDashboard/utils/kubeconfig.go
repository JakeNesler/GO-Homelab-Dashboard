package utils

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// LoadKubeClientsFromSingleFile loads Kubernetes clients for all contexts in the kubeconfig file
func LoadKubeClientsFromSingleFile() ([]*kubernetes.Clientset, error) {
	var clients []*kubernetes.Clientset

	kubeconfigPath := "./configs/kubeconfigs/config" // Path to your kubeconfig file

	// Load the kubeconfig file
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	// Iterate over contexts and create a client for each
	for contextName := range config.Contexts {
		clientConfig, err := clientcmd.NewNonInteractiveClientConfig(*config, contextName, nil, nil).ClientConfig()
		if err != nil {
			fmt.Printf("failed to create client for context %s: %v\n", contextName, err)
			continue
		}

		clientset, err := kubernetes.NewForConfig(clientConfig)
		if err != nil {
			fmt.Printf("failed to create Kubernetes client for context %s: %v\n", contextName, err)
			continue
		}

		clients = append(clients, clientset)
	}

	return clients, nil
}
