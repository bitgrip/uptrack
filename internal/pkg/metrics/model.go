package metrics

// Registry is a datastore to collect metrics
type Registry interface {
	// General
	IncExecution(job string)
	IncCanConnect(job string)
	IncCanNotConnect(job string)
	// SSL Check
	SetSSLDaysLeft(job string, daysLeft float64)
	// HTTP Check
	SetConnectTime(job string, millis int64)
	SetTTFB(job string, millis int64)
	SetRequestTime(job string, millis int64)
	SetBytesReceived(job string, bytes int64)

	//DNS lookup check
	SetIpsRatio(job string, ratio float64)
}
