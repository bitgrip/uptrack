package job

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Descriptor represents a JOB document
type Descriptor struct {
	Name    string            `json:"project,omitempty" yaml:"project,omitempty"`
	UpJobs  map[string]*UpJob `json:"up_jobs,omitempty" yaml:"up_jobs,omitempty"`
	DNSJobs map[string]DnsJob `json:"dns_jobs,omitempty" yaml:"dns_jobs,omitempty"`
}

func DescriptorFromFile(path string) (Descriptor, error) {
	d := Descriptor{}
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
	Name                string
	Host                string            `json:"host" yaml:"host"`
	URL                 string            `json:"url" yaml:"url"`
	Method              Method            `json:"method,omitempty" yaml:"method,omitempty"`
	Expected            int               `json:"expected_code,omitempty" yaml:"expected_code,omitempty"`
	Headers             map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	PlainBody           string            `json:"plain_body,omitempty" yaml:"plain_body,omitempty"`
	Base64Body          string            `json:"base64_body,omitempty" yaml:"base64_body,omitempty"`
	CheckSSL            bool              `json:"check_ssl,omitempty" yaml:"check_ssl,omitempty"`
	CustomTags          map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
	ContentMatch        string            `json:"content_match,omitempty" yaml:"content_match,omitempty"`
	ReverseContentMatch bool              `json:"reverse_content_match,omitempty" yaml:"reverse_content_match,omitempty"`
	Oauth               Oauth             `json:"Oauth,omitempty" yaml:"Oauth,omitempty"`
	OauthServerResponse OauthResponse
}

type Oauth struct {
	AuthUrl string            `json:"auth_url,omitempty" yaml:"auth_url,omitempty"`
	Params  map[string]string `json:"params,omitempty" yaml:"params,omitempty"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type OauthRequest struct {
	GrantType    string `json:"grant_type,omitempty" yaml:"grant_type,omitempty"`
	ClientId     string `json:"client_id,omitempty" yaml:"client_id,omitempty"`
	ClientSecret int    `json:"client_secret,omitempty" yaml:"client_secret,omitempty"`
	Username     string `json:"username,omitempty" yaml:"username,omitempty"`
	Password     string `json:"password,omitempty" yaml:"password,omitempty"`
}

type OauthResponse struct {
	AccessToken  string `json:"access_token,omitempty" yaml:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty" yaml:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty" yaml:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty" yaml:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty" yaml:"scope,omitempty"`
	ExpiresAt    time.Time
	Refresh      bool
}

type UpJobs map[string]UpJob

//Defining default values for unmarshalling UpJob
type tmpUpJob UpJob

func (j *UpJob) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := &tmpUpJob{
		Method:              GET,
		Expected:            200,
		ContentMatch:        "",
		ReverseContentMatch: false,
		CheckSSL:            true,
	}
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
	Name       string
	Host       string            `json:"host" yaml:"host"`
	FQDN       string            `json:"fqdn" yaml:"fqdn"`
	IPs        []string          `json:"ips" yaml:"ips"`
	CustomTags map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
}
