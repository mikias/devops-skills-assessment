package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClientset creates a new Kubernetes clientset from the default kubeconfig file
func NewClientset() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/mike/.kube/config")
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %v", err)
	}
	return kubernetes.NewForConfig(config)
}
