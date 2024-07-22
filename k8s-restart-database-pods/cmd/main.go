package main

import (
	"k8s-restart-database-pods/pkg/k8s" // Adjust to your module's import path
	"log"
)

func main() {
	clientset, err := k8s.NewClientset()
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	log.Println("Starting to restart deployments...")
	if err := k8s.RestartDatabaseDeployments(clientset); err != nil {
		log.Fatalf("Failed to restart deployments: %v", err)
	}
	log.Println("Successfully restarted deployments.")
}
