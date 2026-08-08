package analyzer

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ClusterReader is the seam between a scan and a cluster: it provides the
// Snapshot of cluster state that detectors run against.
type ClusterReader interface {
	// Snapshot lists every resource kind once. Kinds that fail to list are
	// recorded in Snapshot.Unavailable instead of failing the whole scan.
	Snapshot(ctx context.Context) Snapshot
}

// NewClusterReader returns a ClusterReader backed by a live clientset.
func NewClusterReader(clientset *kubernetes.Clientset) ClusterReader {
	return kubeCluster{clientset: clientset}
}

// kubeCluster adapts *kubernetes.Clientset to the ClusterReader seam.
type kubeCluster struct {
	clientset *kubernetes.Clientset
}

func (c kubeCluster) Snapshot(ctx context.Context) Snapshot {
	s := Snapshot{Unavailable: make(map[Kind]error)}

	if deps, err := c.clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindDeployments] = err
	} else {
		s.Deployments = deps.Items
	}

	if pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindPods] = err
	} else {
		s.Pods = pods.Items
	}

	if pvcs, err := c.clientset.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindPVCs] = err
	} else {
		s.PVCs = pvcs.Items
	}

	if svcs, err := c.clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindServices] = err
	} else {
		s.Services = svcs.Items
	}

	if slices, err := c.clientset.DiscoveryV1().EndpointSlices("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindEndpointSlices] = err
	} else {
		s.EndpointSlices = slices.Items
	}

	return s
}
