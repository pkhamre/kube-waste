package analyzer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const CostPerGB = 0.10

// DetectPVCWaste finds PersistentVolumeClaims no running Pod mounts.
func DetectPVCWaste(s Snapshot) []WasteItem {
	var waste []WasteItem

	// Build a map of mounted PVCs (key: namespace/name)
	mountedPVCs := make(map[string]bool)
	for _, pod := range s.Pods {
		for _, vol := range pod.Spec.Volumes {
			if vol.PersistentVolumeClaim != nil {
				key := fmt.Sprintf("%s/%s", pod.Namespace, vol.PersistentVolumeClaim.ClaimName)
				mountedPVCs[key] = true
			}
		}
	}

	for _, pvc := range s.PVCs {
		key := fmt.Sprintf("%s/%s", pvc.Namespace, pvc.Name)
		if mountedPVCs[key] {
			continue
		}

		// client-go handles the "500Mi" vs "1Gi" parsing for us automatically
		qty := pvc.Spec.Resources.Requests["storage"]
		gb := convertToGB(qty)
		cost := gb * CostPerGB

		// Only Bound or Pending PVCs count
		if pvc.Status.Phase == corev1.ClaimBound || pvc.Status.Phase == corev1.ClaimPending {
			waste = append(waste, WasteItem{
				Type:      WastePVC,
				Name:      pvc.Name,
				Namespace: pvc.Namespace,
				Details:   fmt.Sprintf("%s (%s)", qty.String(), pvc.Status.Phase),
				Cost:      cost,
			})
		}
	}

	return waste
}

// convertToGB converts a K8s quantity to simple float GiB.
func convertToGB(q resource.Quantity) float64 {
	// Value() returns raw bytes
	bytes := q.Value()
	return float64(bytes) / (1024 * 1024 * 1024)
}
