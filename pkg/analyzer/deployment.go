package analyzer

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DetectZombieDeployments(clientset *kubernetes.Clientset) ([]WasteItem, error) {
	var waste []WasteItem

	// 1. List all Deployments
	deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, dep := range deployments.Items {
		// 2. Check Replicas
		// In Go K8s structs, Replicas is a pointer (*int32) to distinguish between "0" and "unset"
		if dep.Spec.Replicas != nil && *dep.Spec.Replicas == 0 {
			waste = append(waste, WasteItem{
				Type:      WasteDeploy,
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Details:   "Replicas: 0 (Scaled Down)",
				Cost:      0.0, // Operational waste, not direct billable waste (usually)
			})
		}
	}

	return waste, nil
}
