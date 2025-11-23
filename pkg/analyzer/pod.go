package analyzer

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Average cloud costs (Modify as needed)
const (
	CostPerVCPU  = 20.00
	CostPerGBRam = 3.00
)

const DefaultCPU = 0.1   // 100m
const DefaultMem = 0.125 // 128Mi

func DetectOrphanedPods(clientset *kubernetes.Clientset) ([]WasteItem, error) {
	var waste []WasteItem

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		// 1. Check for OwnerReferences (Deployment, DaemonSet, StatefulSet)
		if len(pod.OwnerReferences) == 0 && pod.Status.Phase == "Running" {

			// Calculate Resources
			cpuReq := 0.0
			memReq := 0.0

			for _, container := range pod.Spec.Containers {
				if req, ok := container.Resources.Requests["cpu"]; ok {
					cpuReq += float64(req.MilliValue()) / 1000.0
				}
				if req, ok := container.Resources.Requests["memory"]; ok {
					// Value() returns bytes, convert to GB
					memReq += float64(req.Value()) / (1024 * 1024 * 1024)
				}
			}

			// If no requests are set, assume a small default for estimation purposes
			isEstimated := false
			if cpuReq == 0 {
				cpuReq = DefaultCPU
				isEstimated = true
			}
			if memReq == 0 {
				memReq = DefaultMem
				isEstimated = true
			}

			cost := (cpuReq * CostPerVCPU) + (memReq * CostPerGBRam)
			details := fmt.Sprintf("%.1f vCPU / %.1f GB", cpuReq, memReq)
			if isEstimated {
				details += " (Est)"
			}

			waste = append(waste, WasteItem{
				Type:      WastePod,
				Name:      pod.Name,
				Namespace: pod.Namespace,
				Details:   details,
				Cost:      cost,
			})
		}
	}

	return waste, nil
}
