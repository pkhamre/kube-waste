package analyzer

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

// DetectorError records that a detector was skipped and why.
type DetectorError struct {
	Detector WasteType
	Err      error
}

// ScanResult is the outcome of a scan: the waste found, and any detectors that
// had to be skipped because a kind they need could not be listed.
type ScanResult struct {
	Waste  []WasteItem
	Errors []DetectorError
}

type detector struct {
	kind  WasteType
	needs []Kind
	run   func(Snapshot) []WasteItem
}

// detectors runs in output order.
var detectors = []detector{
	{kind: WastePVC, needs: []Kind{KindPVCs, KindPods}, run: DetectPVCWaste},
	{kind: WasteDeploy, needs: []Kind{KindDeployments}, run: DetectZombieDeployments},
	{kind: WasteService, needs: []Kind{KindServices, KindEndpointSlices}, run: DetectUnusedServices},
	{kind: WastePod, needs: []Kind{KindPods}, run: DetectOrphanedPods},
}

// Scan lists the cluster once into a Snapshot, then runs every detector whose
// kinds are available, collecting partial results and per-detector errors.
func Scan(ctx context.Context, clientset *kubernetes.Clientset) ScanResult {
	s := fetchSnapshot(ctx, clientset)

	var result ScanResult
	for _, d := range detectors {
		skip := false
		for _, k := range d.needs {
			if !s.Available(k) {
				result.Errors = append(result.Errors, DetectorError{Detector: d.kind, Err: s.Unavailable[k]})
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		result.Waste = append(result.Waste, d.run(s)...)
	}
	return result
}
