package analyzer

import (
	"context"
	"fmt" // <--- We need this back for the label selector

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const CostLoadBalancer = 15.00

func DetectUnusedServices(clientset *kubernetes.Clientset) ([]WasteItem, error) {
	var waste []WasteItem

	// 1. List Services
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, svc := range services.Items {
		// Skip Kubernetes internal service
		if svc.Name == "kubernetes" {
			continue
		}

		// 2. Check EndpointSlices (The modern way to find backing pods)
		// We filter by the service name label to find slices belonging to this service
		labelSelector := fmt.Sprintf("kubernetes.io/service-name=%s", svc.Name)
		slices, err := clientset.DiscoveryV1().EndpointSlices(svc.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})

		hasBackends := false
		if err == nil {
			for _, slice := range slices.Items {
				if len(slice.Endpoints) > 0 {
					hasBackends = true
					break
				}
			}
		}

		if !hasBackends {
			cost := 0.0
			details := "No active pods"

			// Check if it's costing money
			if svc.Spec.Type == "LoadBalancer" {
				cost = CostLoadBalancer
				details = "LoadBalancer (Unused)"
			} else if svc.Spec.Type == "NodePort" {
				details = "NodePort (Unused)"
			} else {
				details = "ClusterIP (Unused)"
			}

			waste = append(waste, WasteItem{
				Type:      WasteService,
				Name:      svc.Name,
				Namespace: svc.Namespace,
				Details:   details,
				Cost:      cost,
			})
		}
	}

	return waste, nil
}
