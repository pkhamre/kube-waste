package analyzer

import "fmt"

// Refs is a Snapshot plus the reference relations derived from it: which PVCs
// are mounted by a live Pod, and which Services have live backends. Detectors
// consult these relations instead of re-deriving them from the raw lists.
type Refs struct {
	Snapshot
	mounted  map[string]bool // key: namespace/name of a PVC mounted by a Pod
	backends map[string]bool // key: namespace/name of a Service with a live EndpointSlice
}

// BuildRefs derives the reference relations from a Snapshot in a single pass.
func BuildRefs(s Snapshot) Refs {
	refs := Refs{Snapshot: s, mounted: make(map[string]bool), backends: make(map[string]bool)}

	for _, pod := range s.Pods {
		for _, vol := range pod.Spec.Volumes {
			if vol.PersistentVolumeClaim != nil {
				key := fmt.Sprintf("%s/%s", pod.Namespace, vol.PersistentVolumeClaim.ClaimName)
				refs.mounted[key] = true
			}
		}
	}

	for _, slice := range s.EndpointSlices {
		if len(slice.Endpoints) > 0 {
			key := fmt.Sprintf("%s/%s", slice.Namespace, slice.Labels["kubernetes.io/service-name"])
			refs.backends[key] = true
		}
	}

	return refs
}

// PVCInUse reports whether a PVC in the given namespace is mounted by a Pod.
func (r Refs) PVCInUse(namespace, name string) bool {
	return r.mounted[fmt.Sprintf("%s/%s", namespace, name)]
}

// ServiceInUse reports whether a Service in the given namespace has a live
// EndpointSlice (one with at least one endpoint).
func (r Refs) ServiceInUse(namespace, name string) bool {
	return r.backends[fmt.Sprintf("%s/%s", namespace, name)]
}
