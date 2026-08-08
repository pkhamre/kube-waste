package analyzer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// DetectPVCWaste finds PersistentVolumeClaims no running Pod mounts.
func DetectPVCWaste(s Refs) []WasteItem {
	var waste []WasteItem

	for _, pvc := range s.PVCs {
		if s.PVCInUse(pvc.Namespace, pvc.Name) {
			continue
		}

		// client-go handles the "500Mi" vs "1Gi" parsing for us automatically
		qty := pvc.Spec.Resources.Requests["storage"]
		gb := convertToGB(qty)

		// Only Bound or Pending PVCs count
		if pvc.Status.Phase == corev1.ClaimBound || pvc.Status.Phase == corev1.ClaimPending {
			waste = append(waste, WasteItem{
				Type:      WastePVC,
				Name:      pvc.Name,
				Namespace: pvc.Namespace,
				Details:   fmt.Sprintf("%s (%s)", qty.String(), pvc.Status.Phase),
				Usage:     Usage{StorageGB: gb},
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
