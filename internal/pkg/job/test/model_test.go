package test

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/assert"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"testing"
)

func TestDescriptorUnmarshal(t *testing.T) {
	descriptor, err := job.DescriptorFromFile("./test.yaml")
	assert.Equals(t, nil, err)
	assert.Equals(t, 3, len(descriptor.UpJobs))

	upJob := descriptor.UpJobs["bitgrip_checker"]
	name := upJob.Name
	assert.Equals(t, "bitgrip_checker", name)

	url := upJob.URL
	assert.Equals(t, "https://www.bitgrip.de/kontakt", url)

	ipsCount := len(descriptor.DNSJobs["bitgrip_dns_check"].IPs)
	assert.Equals(t, 4, ipsCount)

	tagsCount := len(upJob.CustomTags)
	assert.Equals(t, 3, tagsCount)

	headers := upJob.Headers
	assert.Equals(t, "Basic ABC==", headers["Authorization"])

	host, _ := upJob.HostString()
	assert.Equals(t, "bitgrip.de", host)

	other1 := descriptor.UpJobs["other1"]

	host, _ = other1.HostString()
	assert.Equals(t, "bla.com", host)

	other2 := descriptor.UpJobs["other2"]

	host, _ = other2.HostString()
	assert.Equals(t, "ci.bitgrip.de", host)

}
