package job

// Descriptor represents a JOB document
type Descriptor struct {
	BaseURL   string     `json:"base_url,omitempty" yaml:"base_url,omitempty"`
	UpChecks  []UpCheck  `json:"up_checks,omitempty" yaml:"up_checks,omitempty"`
	DNSChecks []DNSCheck `json:"dns_checks,omitempty" yaml:"dns_checks,omitempty"`
}

// UpCheck is a check is a HTTP endpoint is up and able to server required method
type UpCheck struct {
	Method     Method   `json:"method,omitempty" yaml:"method,omitempty"`
	URI        string   `json:"uri,omitempty" yaml:"uri,omitempty"`
	Headers    []Header `json:"headers,omitempty" yaml:"headers,omitempty"`
	PlainBody  string   `json:"plain_body,omitempty" yaml:"plain_body,omitempty"`
	Base64Body string   `json:"base64_body,omitempty" yaml:"base64_body,omitempty"`
	CheckSSL   bool     `json:"check_ssl,omitempty" yaml:"check_ssl,omitempty"`
}

// Method is a HTTP method
type Method string

const (
	// GET is used to receive an entity from a HTTP endpoint
	GET Method = "GET"
	// POST is used to post data to a HTTP endpoint
	POST Method = "POST"
	// PUT is used to send an entity to a HTTP endpoint
	PUT Method = "PUT"
	// DELETE is used to delete an entity to a HTTP endpoint
	DELETE Method = "DELETE"
	// PATCH is used to update an entity at a HTTP endpoint
	PATCH Method = "PATCH"
	// OPTIONS is used to verify CORS policies
	OPTIONS Method = "OPTIONS"
	// HEAD is used to receive headers of an entity from a HTTP endpoint
	HEAD Method = "HEAD"
)

// Header represent HTTP headers to be used in a check
type Header map[string][]string

// DNSCheck is verifying if a fqdn is looked up to the expected set of IPs
type DNSCheck struct {
	FQDN     string   `json:"fqdn,omitempty" yaml:"fqdn,omitempty"`
	IPs      []string `json:"ips,omitempty" yaml:"ips,omitempty"`
	CheckSSL bool     `json:"check_ssl,omitempty" yaml:"check_ssl,omitempty"`
}
