package job

// Descriptor represents a JOB document
type Descriptor struct {
	BaseURL string           `json:"base_url,omitempty" yaml:"base_url,omitempty"`
	UpJobs  map[string]UpJob `json:"up_jobs,omitempty" yaml:"up_jobs,omitempty"`
	DNSJobs []DnsJob         `json:"dns_jobs,omitempty" yaml:"dns_jobs,omitempty"`
}

// UpJob is a check if a HTTP endpoint is up and able to serve required method
type UpJob struct {
	Name       string
	Method     Method `json:"method,omitempty" yaml:"method,omitempty"`
	Expected   int    `json:"expected_code,omitempty" yaml:"expected_code,omitempty"`
	URLSuffix  string `json:"url_suffix,omitempty" yaml:"url_suffix,omitempty"`
	URL        string
	Headers    []Header `json:"headers,omitempty" yaml:"headers,omitempty"`
	PlainBody  string   `json:"plain_body,omitempty" yaml:"plain_body,omitempty"`
	Base64Body string   `json:"base64_body,omitempty" yaml:"base64_body,omitempty"`
	CheckSSL   bool     `json:"check_ssl,omitempty" yaml:"check_ssl,omitempty"`
}

type UpJobs map[string]UpJob

//Defining default values for unmarshalling UpJob
func (j *UpJob) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tmpUpJob UpJob
	test := &tmpUpJob{
		Method:   GET,
		Expected: 200,
		CheckSSL: true}
	unmarshal(&test)
	*j = UpJob(*test)
	return nil
}

// Method is a HTTP method
type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	PATCH   Method = "PATCH"
	OPTIONS Method = "OPTIONS"
	HEAD    Method = "HEAD"
)

// Header represent HTTP headers to be used in a check
type Header map[string][]string

// DnsJob is verifying if a fqdn is looked up to the expected set of IPs
type DnsJob struct {
	FQDN     string   `json:"fqdn,omitempty" yaml:"fqdn,omitempty"`
	IPs      []string `json:"ips,omitempty" yaml:"ips,omitempty"`
	CheckSSL bool     `json:"check_ssl,omitempty" yaml:"check_ssl,omitempty"`
}
