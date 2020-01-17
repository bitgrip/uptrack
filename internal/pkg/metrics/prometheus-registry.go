package metrics

// prometheusRegistry is a wrapper to forward Registry actions
// to a collection of Registries
type prometheusRegistry struct {
}

func NewPrometheusRegistry(listenOn string) Registry {
	return &prometheusRegistry{}
}

func (r *prometheusRegistry) IncExecution(job string) {
}

func (r *prometheusRegistry) IncCanConnect(job string, uri string) {
}

func (r *prometheusRegistry) IncCanNotConnect(job string, uri string) {
}

func (r *prometheusRegistry) SetSSLDaysLeft(job string, uri string, daysLeft int64) {
}

func (r *prometheusRegistry) SetConnectTime(job string, uri string, millis int64) {
}

func (r *prometheusRegistry) SetTTFB(job string, uri string, millis int64) {
}

func (r *prometheusRegistry) SetRequestTime(job string, uri string, millis int64) {
}

func (r *prometheusRegistry) SetBytesReceived(job string, uri string, bytes int64) {
}
