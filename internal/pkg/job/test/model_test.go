package test

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/assert"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"testing"
)

func TestDescriptorUnmarshal(t *testing.T) {
	descriptor, err := job.DescriptorFromFile("./test.yaml")
	assert.Equals(t, nil, err)
	assert.Equals(t, 2, len(descriptor.UpJobs))
	name := descriptor.UpJobs["bitgrip_checker"].Name
	assert.Equals(t, "bitgrip_checker", name)
	url := descriptor.UpJobs["bitgrip_checker"].URL
	assert.Equals(t, "https://www.bitgrip.de/kontakt", url)

}
