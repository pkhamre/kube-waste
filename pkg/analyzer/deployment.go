package analyzer

// DetectZombieDeployments finds Deployments scaled down to zero replicas.
func DetectZombieDeployments(s Snapshot) []WasteItem {
	var waste []WasteItem

	for _, dep := range s.Deployments {
		// Replicas is a pointer (*int32) to distinguish between "0" and "unset"
		if dep.Spec.Replicas != nil && *dep.Spec.Replicas == 0 {
			waste = append(waste, WasteItem{
				Type:      WasteDeploy,
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Details:   "Replicas: 0 (Scaled Down)",
				Cost:      0.0, // Operational waste, not direct billable waste (usually)
			})
		}
	}

	return waste
}
