package k8s

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RestartDatabaseDeployments restarts all deployments with "database" in their pod names.
func RestartDatabaseDeployments(clientset *kubernetes.Clientset) error {
	// List all pods with the word "database" in their name
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}

	deployments := make(map[string]string) // maps namespace to deployment name

	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "database") {
			for _, owner := range pod.OwnerReferences {
				if owner.Kind == "ReplicaSet" {
					// Find the parent Deployment of the ReplicaSet
					rs, err := clientset.AppsV1().ReplicaSets(pod.Namespace).Get(context.TODO(), owner.Name, metav1.GetOptions{})
					if err != nil {
						return fmt.Errorf("failed to get ReplicaSet %s: %v", owner.Name, err)
					}
					if rs.OwnerReferences != nil {
						for _, depOwner := range rs.OwnerReferences {
							if depOwner.Kind == "Deployment" {
								deployments[pod.Namespace] = depOwner.Name
							}
						}
					}
				}
			}
		}
	}

	// Restart deployments
	for namespace, name := range deployments {
		log.Printf("Restarting deployment: %s in namespace: %s\n", name, namespace)
		if err := rolloutRestartDeployment(clientset, namespace, name); err != nil {
			log.Printf("Failed to restart deployment %s: %v", name, err)
		}
	}

	return nil
}

func rolloutRestartDeployment(clientset *kubernetes.Clientset, namespace, name string) error {
	// Get the deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment %s: %v", name, err)
	}

	// Create a new annotation to trigger the rollout restart
	annotations := deployment.Spec.Template.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	deployment.Spec.Template.Annotations = annotations

	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	return err
}
