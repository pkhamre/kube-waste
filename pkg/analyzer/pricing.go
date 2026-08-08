package analyzer

// Pricing is the monthly cost model: the rate charged per unit of usage.
type Pricing struct {
	VCPU         float64 // per vCPU per month
	MemGB        float64 // per GB of RAM per month
	StorageGB    float64 // per GB of storage per month
	LoadBalancer float64 // per LoadBalancer service per month
}

// DefaultPricing returns the standard average cloud-provider rates.
func DefaultPricing() Pricing {
	return Pricing{VCPU: 20.00, MemGB: 3.00, StorageGB: 0.10, LoadBalancer: 15.00}
}

// Price returns the estimated monthly cost of a waste item of type t with the
// given resource usage.
func (p Pricing) Price(t WasteType, u Usage) float64 {
	switch t {
	case WastePod:
		return (u.CPU * p.VCPU) + (u.MemGB * p.MemGB)
	case WastePVC:
		return u.StorageGB * p.StorageGB
	case WasteService:
		if u.LoadBalancer {
			return p.LoadBalancer
		}
		return 0
	default:
		// WasteDeploy is operational waste, not direct billable waste
		return 0
	}
}
