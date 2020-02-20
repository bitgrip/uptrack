package job

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

// Descriptor represents a JOB document
type Descriptor struct {
	Name              string            `json:"project,omitempty" yaml:"project,omitempty"`
	UpJobs            map[string]UpJob  `json:"up_jobs,omitempty" yaml:"up_jobs,omitempty"`
	DNSJobs           map[string]DnsJob `json:"dns_jobs,omitempty" yaml:"dns_jobs,omitempty"`
	DatadogEnabled    bool              `json:"datadog_enabled,omitempty" yaml:"datadog_enabled,omitempty"`
	PrometheusEnabled bool              `json:"prometheus_enabled,omitempty" yaml:"prometheus_enabled,omitempty"`
}

func DescriptorFromFile(path string) (Descriptor, error) {
	d := Descriptor{
		DatadogEnabled:    true,
		PrometheusEnabled: true,
	}
	data, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(data, &d)

	for name, upJob := range d.UpJobs {
		upJob.Name = name
		d.UpJobs[name] = upJob
	}

	for name, dnsJob := range d.DNSJobs {
		dnsJob.Name = name
		d.DNSJobs[name] = dnsJob
	}
	return d, err
}

// UpJob is a check if a HTTP endpoint is up and able to serve required method
type UpJob struct {
	Name       string
	Host       string      `json:"host" yaml:"host"`
	URL        string      `json:"url" yaml:"url"`
	Method     Method      `json:"method,omitempty" yaml:"method,omitempty"`
	Expected   int         `json:"expected_code,omitempty" yaml:"expected_code,omitempty"`
	Header     http.Header `json:"header,omitempty" yaml:"header,omitempty"`
	PlainBody  string      `json:"plain_body,omitempty" yaml:"plain_body,omitempty"`
	Base64Body string      `json:"base64_body,omitempty" yaml:"base64_body,omitempty"`
	CheckSSL   bool        `json:"check_ssl,omitempty" yaml:"check_ssl,omitempty"`
}

type UpJobs map[string]UpJob

//Defining default values for unmarshalling UpJob
type tmpUpJob UpJob

func (j *UpJob) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := &tmpUpJob{
		Method:   GET,
		Expected: 200,
		CheckSSL: true}
	unmarshal(tmp)
	*j = UpJob(*tmp)
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

// DnsJob is verifying if a fqdn is looked up to the expected set of IPs
type DnsJob struct {
	Name string
	Host string   `json:"host" yaml:"host"`
	FQDN string   `json:"fqdn" yaml:"fqdn"`
	IPs  []string `json:"ips" yaml:"ips"`
}
