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

func (r *combinedRegistry) IncExecution(job string) {
	for _, registry := range r.registries {
		registry.IncExecution(job)
	}
}

func (r *combinedRegistry) IncCanConnect(job string) {
	for _, registry := range r.registries {
		registry.IncCanConnect(job)
	}
}

func (r *combinedRegistry) IncCanNotConnect(job string) {
	for _, registry := range r.registries {
		registry.IncCanNotConnect(job)
	}
}

func (r *combinedRegistry) SetSSLDaysLeft(job string, daysLeft float64) {
	for _, registry := range r.registries {
		registry.SetSSLDaysLeft(job, daysLeft)
	}
}

func (r *combinedRegistry) SetConnectTime(job string, millis int64) {
	for _, registry := range r.registries {
		registry.SetConnectTime(job, millis)
	}
}

func (r *combinedRegistry) SetTTFB(job string, millis int64) {
	for _, registry := range r.registries {
		registry.SetTTFB(job, millis)
	}
}

func (r *combinedRegistry) SetRequestTime(job string, millis int64) {
	for _, registry := range r.registries {
		registry.SetRequestTime(job, millis)
	}
}

func (r *combinedRegistry) SetBytesReceived(job string, bytes int64) {
	for _, registry := range r.registries {
		registry.SetBytesReceived(job, bytes)
	}
}
