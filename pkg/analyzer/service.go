package analyzer

import (
	corev1 "k8s.io/api/core/v1"
)

const CostLoadBalancer = 15.00

// DetectUnusedServices finds Services with no backing Pods, filtering by
// EndpointSlices.
func DetectUnusedServices(s Snapshot) []WasteItem {
	var waste []WasteItem

	for _, svc := range s.Services {
		// Skip Kubernetes internal service
		if svc.Name == "kubernetes" {
			continue
		}

		// Find endpoint slices belonging to this service
		hasBackends := false
		for _, slice := range s.EndpointSlices {
			if slice.Namespace == svc.Namespace &&
				slice.Labels["kubernetes.io/service-name"] == svc.Name &&
				len(slice.Endpoints) > 0 {
				hasBackends = true
				break
			}
		}

		if hasBackends {
			continue
		}

		cost := 0.0
		details := "No active pods"

		// Check if it's costing money
		switch svc.Spec.Type {
		case corev1.ServiceTypeLoadBalancer:
			cost = CostLoadBalancer
			details = "LoadBalancer (Unused)"
		case corev1.ServiceTypeNodePort:
			details = "NodePort (Unused)"
		default:
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

	return waste
}
