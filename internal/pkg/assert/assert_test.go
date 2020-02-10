// Copyright © 2018 Bitgrip <berlin@bitgrip.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package assert

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/ctl"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"testing"
)

func TestDescriptorUnmarshal(t *testing.T) {
	descriptor, err := job.DescriptorFromFile("./test.yaml")
	Equals(t, nil, err)
	Equals(t, 2, len(descriptor.UpJobs))
	name := descriptor.UpJobs["bitgrip_checker"].Name
	Equals(t, "bitgrip_checker", name)
	url := descriptor.UpJobs["bitgrip_checker"].URL
	Equals(t, "https://www.bitgrip.de/kontakt", url)

}
func TestIntesectArrays(t *testing.T) {
	arr1 := []string{"A", "B", "C", "D"}
	arr2 := []string{"1", "B", "3", "D"}
	intersecting := ctl.GetIntersecting(arr1, arr2)
	Equals(t, []string{"B", "D"}, intersecting)
}
