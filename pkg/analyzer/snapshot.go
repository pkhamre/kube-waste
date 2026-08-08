package analyzer

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
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
// fetched once via ClusterReader and passed to every detector.
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
