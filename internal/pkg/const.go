package cons

/// GENERAL constants, used for tagging

const (
	JobName     string = "job"
	CheckName   string = "check"
	UrlString   string = "url"
	ReqMethod   string = "request_method"
	Host               = "host"
	FQDN        string = "FQDN"
	ProjectName        = "project"
)

//DataDog metric suffixes
const (
	DDCanConnect    string = "connection.successful"
	DDCannotConnect string = "connection.failed"
	DDSSLDaysLeft   string = "ssl_days_left"
	DDConnectTime   string = "connection.time"
	DDTTFB          string = "TTFB"
	DDRequestTime   string = "request_time"
	DDBytesReceived string = "bytes_received"
	DDFoundIps      string = "found_ips_ratio"
)

//Prometheus metric suffixes
const (
	PromCanConnect    string = "connection_successful"
	PromCannotConnect string = "connection_failed"
	PromSSLDaysLeft   string = "ssl_days_left"
	PromConnectTime   string = "connection_time"
	PromTTFB          string = "TTFB"
	PromRequestTime   string = "request_time"
	PromBytesReceived string = "bytes_received"
	PromFoundIps      string = "found_ips_ratio"

	//Prefixes for metric keys
	PromNamespace          string = "uptrack"
	PromNameUpcheckCounter string = "upcheck_counter"
	PromNameUpCheckGauge   string = "upcheck_gauge"

	PromNameDnsCheckCounter string = "upcheck_counter"
	PromNameDnsCheckGauge   string = "dnscheck_gauge"
)
