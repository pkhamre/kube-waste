package analyzer

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/api/resource"
)

const CostPerGB = 0.10

func DetectPVCWaste(clientset *kubernetes.Clientset) ([]WasteItem, error) {
	var waste []WasteItem

	// 1. Get all PVCs
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 2. Get all Pods to check mounts
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 3. Create a map of mounted PVCs
	mountedPVCs := make(map[string]bool)
	for _, pod := range pods.Items {
		for _, vol := range pod.Spec.Volumes {
			if vol.PersistentVolumeClaim != nil {
				// Key format: namespace/name
				key := fmt.Sprintf("%s/%s", pod.Namespace, vol.PersistentVolumeClaim.ClaimName)
				mountedPVCs[key] = true
			}
		}
	}

	// 4. Find the orphans
	for _, pvc := range pvcs.Items {
		key := fmt.Sprintf("%s/%s", pvc.Namespace, pvc.Name)
		if !mountedPVCs[key] {
			// Calculate Size
			// client-go handles the "500Mi" vs "1Gi" parsing for us automatically!
			qty := pvc.Spec.Resources.Requests["storage"] 
			gb := convertToGB(qty)
			cost := gb * CostPerGB

			// Filter Logic (Show Bound or Pending)
			if pvc.Status.Phase == "Bound" || pvc.Status.Phase == "Pending" {
				waste = append(waste, WasteItem{
					Type:      WastePVC,
					Name:      pvc.Name,
					Namespace: pvc.Namespace,
					Details:   fmt.Sprintf("%s (%s)", qty.String(), pvc.Status.Phase),
					Cost:      cost,
				})
			}
		}
	}

	return waste, nil
}

// Helper to convert weird K8s quantities to simple Float GBs
func convertToGB(q resource.Quantity) float64 {
	// Value() returns raw bytes. 
	bytes := q.Value()
	return float64(bytes) / (1024 * 1024 * 1024)
}
