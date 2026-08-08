package analyzer

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Kind identifies a cluster resource kind that a Snapshot can hold.
type Kind string

const (
	KindDeployments    Kind = "deployments"
	KindPods           Kind = "pods"
	KindPVCs           Kind = "persistentvolumeclaims"
	KindServices       Kind = "services"
	KindEndpointSlices Kind = "endpointslices"
)

// Snapshot is the immutable cluster state a single scan runs against. It is
// fetched once by Scan and passed to every detector.
type Snapshot struct {
	Deployments    []appsv1.Deployment
	Pods           []corev1.Pod
	PVCs           []corev1.PersistentVolumeClaim
	Services       []corev1.Service
	EndpointSlices []discoveryv1.EndpointSlice
	// Unavailable records kinds that failed to list. Detectors that need an
	// unavailable kind must be skipped rather than run against empty data.
	Unavailable map[Kind]error
}

// Available reports whether the kind was listed successfully.
func (s *Snapshot) Available(k Kind) bool {
	_, ok := s.Unavailable[k]
	return !ok
}

// fetchSnapshot lists every resource kind once. Kinds that fail to list are
// recorded in Unavailable instead of aborting the whole scan.
func fetchSnapshot(ctx context.Context, clientset *kubernetes.Clientset) Snapshot {
	s := Snapshot{Unavailable: make(map[Kind]error)}

	if deps, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindDeployments] = err
	} else {
		s.Deployments = deps.Items
	}

	if pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindPods] = err
	} else {
		s.Pods = pods.Items
	}

	if pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindPVCs] = err
	} else {
		s.PVCs = pvcs.Items
	}

	if svcs, err := clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindServices] = err
	} else {
		s.Services = svcs.Items
	}

	if slices, err := clientset.DiscoveryV1().EndpointSlices("").List(ctx, metav1.ListOptions{}); err != nil {
		s.Unavailable[KindEndpointSlices] = err
	} else {
		s.EndpointSlices = slices.Items
	}

	return s
}
