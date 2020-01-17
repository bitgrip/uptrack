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

func (r *combinedRegistry) IncCanConnect(job string, uri string) {
	for _, registry := range r.registries {
		registry.IncCanConnect(job, uri)
	}
}

func (r *combinedRegistry) IncCanNotConnect(job string, uri string) {
	for _, registry := range r.registries {
		registry.IncCanNotConnect(job, uri)
	}
}

func (r *combinedRegistry) SetSSLDaysLeft(job string, uri string, daysLeft int64) {
	for _, registry := range r.registries {
		registry.SetSSLDaysLeft(job, uri, daysLeft)
	}
}

func (r *combinedRegistry) SetConnectTime(job string, uri string, millis int64) {
	for _, registry := range r.registries {
		registry.SetConnectTime(job, uri, millis)
	}
}

func (r *combinedRegistry) SetTTFB(job string, uri string, millis int64) {
	for _, registry := range r.registries {
		registry.SetTTFB(job, uri, millis)
	}
}

func (r *combinedRegistry) SetRequestTime(job string, uri string, millis int64) {
	for _, registry := range r.registries {
		registry.SetRequestTime(job, uri, millis)
	}
}

func (r *combinedRegistry) SetBytesReceived(job string, uri string, bytes int64) {
	for _, registry := range r.registries {
		registry.SetBytesReceived(job, uri, bytes)
	}
}
