package metrics

// combinedRegistry is a wrapper to forward Registry actions
// to a collection of Registries
type combinedRegistry struct {
	registries []Registry
}

func NewCombinedRegistry(registries ...Registry) Registry {
	return &combinedRegistry{
		registries: registries,
	}
}

func (r *combinedRegistry) CanConnect(job string) {
	for _, registry := range r.registries {
		registry.CanConnect(job)
	}
}

func (r *combinedRegistry) CanNotConnect(job string) {
	for _, registry := range r.registries {
		registry.CanNotConnect(job)
	}
}

func (r *combinedRegistry) SetSSLDaysLeft(job string, daysLeft float64) {
	for _, registry := range r.registries {
		registry.SetSSLDaysLeft(job, daysLeft)
	}
}

func (r *combinedRegistry) SetConnectTime(job string, millis float64) {
	for _, registry := range r.registries {
		registry.SetConnectTime(job, millis)
	}
}

func (r *combinedRegistry) SetTTFB(job string, millis float64) {
	for _, registry := range r.registries {
		registry.SetTTFB(job, millis)
	}
}

func (r *combinedRegistry) SetRequestTime(job string, millis float64) {
	for _, registry := range r.registries {
		registry.SetRequestTime(job, millis)
	}
}

func (r *combinedRegistry) SetBytesReceived(job string, bytes float64) {
	for _, registry := range r.registries {
		registry.SetBytesReceived(job, bytes)
	}
}

func (r *combinedRegistry) SetIpsRatio(job string, ratio float64) {
	for _, registry := range r.registries {
		registry.SetIpsRatio(job, ratio)
	}
}
