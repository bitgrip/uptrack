package metrics

import "regexp"

// Registry is a datastore to collect metrics
type Registry interface {
	// General
	IncCanConnect(job string)
	IncCanNotConnect(job string)
	// SSL Check
	SetSSLDaysLeft(job string, daysLeft float64)
	// HTTP Check
	SetConnectTime(job string, millis float64)
	SetTTFB(job string, millis float64)
	SetRequestTime(job string, millis float64)
	SetBytesReceived(job string, bytes float64)

	//DNS lookup check
	SetIpsRatio(job string, ratio float64)
}

func replaceAll(str string, pattern string) string {
	r, _ := regexp.Compile(pattern)
	str = string(r.ReplaceAll([]byte(str), []byte("_")))
	return str
}
