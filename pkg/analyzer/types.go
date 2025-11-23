package analyzer

type WasteType string

const (
	WastePVC      WasteType = "Unused PVC"
	WasteService  WasteType = "Unused Service"
	WastePod      WasteType = "Orphaned Pod"
	WasteDeploy   WasteType = "Zombie Deployment"
)

type WasteItem struct {
	Type      WasteType
	Name      string
	Namespace string
	Details   string
	Cost      float64
}
