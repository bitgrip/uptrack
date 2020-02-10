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

func (r *datadogRegistry) IncCanConnect(job string) {
}

func (r *datadogRegistry) IncCanNotConnect(job string) {
}

func (r *datadogRegistry) SetSSLDaysLeft(job string, daysLeft float64) {
}

func (r *datadogRegistry) SetConnectTime(job string, millis int64) {
}

func (r *datadogRegistry) SetTTFB(job string, millis int64) {
}

func (r *datadogRegistry) SetRequestTime(job string, millis int64) {
}

func (r *datadogRegistry) SetBytesReceived(job string, bytes int64) {
}
func (r *datadogRegistry) SetIpsRatio(job string, ratio float64) {

}
