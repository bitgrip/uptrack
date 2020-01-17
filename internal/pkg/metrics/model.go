package metrics

// Registry is a datastore to collect metrics
type Registry interface {
	// General
	IncExecution(job string)
	IncCanConnect(job string, uri string)
	IncCanNotConnect(job string, uri string)
	// SSL Check
	SetSSLDaysLeft(job string, uri string, daysLeft int64)
	// HTTP Check
	SetConnectTime(job string, uri string, millis int64)
	SetTTFB(job string, uri string, millis int64)
	SetRequestTime(job string, uri string, millis int64)
	SetBytesReceived(job string, uri string, bytes int64)
}
