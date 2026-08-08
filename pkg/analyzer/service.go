package analyzer

import (
	corev1 "k8s.io/api/core/v1"
)

// DetectUnusedServices finds Services with no backing Pods, filtering by
// EndpointSlices.
func DetectUnusedServices(s Refs) []WasteItem {
	var waste []WasteItem

	for _, svc := range s.Services {
		// Skip Kubernetes internal service
		if svc.Name == "kubernetes" {
			continue
		}

		if s.ServiceInUse(svc.Namespace, svc.Name) {
			continue
		}

		details := "No active pods"
		switch svc.Spec.Type {
		case corev1.ServiceTypeLoadBalancer:
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
			Usage:     Usage{LoadBalancer: svc.Spec.Type == corev1.ServiceTypeLoadBalancer},
		})
	}

	return waste
}
