package metrics

// datadogRegistry is a wrapper to forward Registry actions
// to a collection of Registries
type datadogRegistry struct {
}

func NewDatadogRegistry(datadogCredentials string) Registry {
	return &datadogRegistry{}
}

func (r *datadogRegistry) IncExecution(job string) {
}

func (r *datadogRegistry) IncCanConnect(job string, uri string) {
}

func (r *datadogRegistry) IncCanNotConnect(job string, uri string) {
}

func (r *datadogRegistry) SetSSLDaysLeft(job string, uri string, daysLeft int64) {
}

func (r *datadogRegistry) SetConnectTime(job string, uri string, millis int64) {
}

func (r *datadogRegistry) SetTTFB(job string, uri string, millis int64) {
}

func (r *datadogRegistry) SetRequestTime(job string, uri string, millis int64) {
}

func (r *datadogRegistry) SetBytesReceived(job string, uri string, bytes int64) {
}
