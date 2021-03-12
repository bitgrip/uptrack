package job

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"reflect"
)

// Descriptor represents a JOB document
type Descriptor struct {
	Name    string             `json:"project,omitempty" yaml:"project,omitempty"`
	UpJobs  map[string]*UpJob  `json:"up_jobs,omitempty" yaml:"up_jobs,omitempty"`
	DNSJobs map[string]*DnsJob `json:"dns_jobs,omitempty" yaml:"dns_jobs,omitempty"`
}

// UpJob is a check if a HTTP endpoint is up and able to serve required method
type UpJob struct {
	Name                string
	Host                string            `json:"host" yaml:"host"`
	URL                 string            `json:"url" yaml:"url"`
	Method              string            `json:"method,omitempty" yaml:"method,omitempty"`
	Expected            int               `json:"expected_code,omitempty" yaml:"expected_code,omitempty"`
	Headers             map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	PlainBody           string            `json:"plain_body,omitempty" yaml:"plain_body,omitempty"`
	Base64Body          string            `json:"base64_body,omitempty" yaml:"base64_body,omitempty"`
	CheckSSL            bool              `json:"check_ssl,omitempty" yaml:"check_ssl,omitempty"`
	CustomTags          map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
	ContentMatch        string            `json:"content_match,omitempty" yaml:"content_match,omitempty"`
	ReverseContentMatch bool              `json:"reverse_content_match,omitempty" yaml:"reverse_content_match,omitempty"`
	Oauth               Oauth             `json:"Oauth,omitempty" yaml:"Oauth,omitempty"`
	Context             Context
	RequestInterceptors []HttpRequestInterceptor
}

func DescriptorFromFile(path string) (Descriptor, error) {
	d := Descriptor{}
	data, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(data, &d)

	for name, upJob := range d.UpJobs {
		upJob.Name = name
		//Check, if upJob implemets Job interface, and call init(), in case it does
		value := reflect.ValueOf(upJob)
		if job, yes := value.Interface().(Job); yes {
			job.init()
		}
	}

	for name, dnsJob := range d.DNSJobs {
		dnsJob.Name = name
	}
	return d, err
}

type Job interface {
	init()
}

type Context struct {
	OauthResponse OauthResponse
}

//Defining default values for unmarshalling UpJob
type tmpUpJob UpJob

func (job *UpJob) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := &tmpUpJob{
		Method:              http.MethodGet,
		Expected:            200,
		ContentMatch:        "",
		ReverseContentMatch: false,
		CheckSSL:            true,
	}
	unmarshal(tmp)
	*job = UpJob(*tmp)
	return nil
}

// DnsJob is verifying if a fqdn is looked up to the expected set of IPs
type DnsJob struct {
	Name       string
	Host       string            `json:"host" yaml:"host"`
	FQDN       string            `json:"fqdn" yaml:"fqdn"`
	IPs        []string          `json:"ips" yaml:"ips"`
	CustomTags map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type HttpRequestInterceptor interface {
	Intercept(req *http.Request, upJob *UpJob) *http.Request
}
