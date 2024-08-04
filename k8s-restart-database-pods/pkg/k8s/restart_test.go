package k8s

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestRestartDatabaseDeployments(t *testing.T) {
	clientset := fake.NewSimpleClientset(
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "database-nginx-1",
				Namespace: "default",
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind: "ReplicaSet",
						Name: "nginx-replicaset",
					},
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "other-pod",
				Namespace: "default",
			},
		},
		&appsv1.ReplicaSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nginx-replicaset",
				Namespace: "default",
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind: "Deployment",
						Name: "database-nginx-deploy",
					},
				},
			},
		},
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "database-nginx-deploy",
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
		},
	)

	if err := RestartDatabaseDeployments(clientset); err != nil {
		t.Fatalf("RestartDatabaseDeployments failed: %v", err)
	}

	// Check if the deployment has been updated
	deployment, err := clientset.AppsV1().Deployments("default").Get(context.TODO(), "database-nginx-deploy", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Failed to get deployment: %v", err)
	}

	restartedAt, ok := deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"]
	if !ok {
		t.Fatalf("Deployment annotation not updated")
	}

	if _, err := time.Parse(time.RFC3339, restartedAt); err != nil {
		t.Fatalf("Failed to parse annotation timestamp: %v", err)
	}
}
