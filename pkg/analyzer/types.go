package analyzer

type WasteType string

const (
	WastePVC      WasteType = "Unused PVC"
	WasteService  WasteType = "Unused Service"
	WastePod      WasteType = "Orphaned Pod"
	WasteDeploy   WasteType = "Zombie Deployment"
)

// Usage is the resource footprint a waste item holds, used by Pricing.
type Usage struct {
	CPU          float64 // vCPU
	MemGB        float64 // GB of RAM
	StorageGB    float64 // GB of storage
	LoadBalancer bool    // backed by a LoadBalancer service
	Estimated    bool    // derived from defaults, not observed requests
}

type WasteItem struct {
	Type      WasteType
	Name      string
	Namespace string
	Details   string
	Usage     Usage
	Cost      float64
}
