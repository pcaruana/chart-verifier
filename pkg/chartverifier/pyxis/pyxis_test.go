/*
 * Copyright 2021 Red Hat
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pyxis

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getImageRegistries(t *testing.T) {

	type testCase struct {
		description string
		repository  string
		registry    string
		message     string
	}

	PassTestCases := []testCase{
		{description: "Test rhscl respository", repository: "rhscl/postgresql-10-rhel7", registry: "registry.access.redhat.com", message: ""},
		{description: "Test rhel6.7 respository", repository: "rhel6.7", registry: "registry.access.redhat.com", message: ""},
		{description: "Test rhel8/nginx-116 respository", repository: "rhel8/nginx-116", registry: "registry.access.redhat.com", message: ""},
		{description: "Test ibm/nginx respository", repository: "ibm/nginx", registry: "non_registry", message: ""},
		{description: "Test turbonomic/zookeeper respository", repository: "turbonomic/zookeeper", registry: "registry.connect.redhat.com", message: ""}}

	for _, tc := range PassTestCases {
		t.Run(tc.description, func(t *testing.T) {
			reg, err := GetImageRegistries(tc.repository)
			require.NoError(t, err)
			require.Equal(t, tc.registry, reg[0])
		})
	}

	FailTestCases := []testCase{
		{description: "Test repository not found", repository: "ndoesnotexist", registry: "registry.hub.docker.com", message: "Respository not found"},
	}

	for _, tc := range FailTestCases {
		t.Run(tc.description, func(t *testing.T) {
			reg, err := GetImageRegistries(tc.repository)
			require.Error(t, err)
			require.Empty(t, reg)
			require.Contains(t, err.Error(), tc.message)
		})
	}

}

func Test_checkImageInRegistry(t *testing.T) {

	type testCase struct {
		description string
		message     string
		imageRef    ImageReference
	}

	PassTestCases := []testCase{
		{description: "Test nginx registry and version found.", message: "", imageRef: ImageReference{Repository: "rhel6.7", Registries: []string{"registry.access.redhat.com"}, Tag: "6.7", Sha: ""}},
		{description: "Test nginx rhel6.7 and version found.", imageRef: ImageReference{Repository: "rhel6.7", Registries: []string{"registry.access.redhat.com"}, Tag: "6.7", Sha: ""}, message: ""},
		{description: "Test rhel8/nginx-116 respository found.", imageRef: ImageReference{Repository: "rhel8/nginx-116", Registries: []string{"registry.access.redhat.com"}, Tag: "1-75", Sha: ""}, message: ""},
		{description: "Test turbonomic/zookeeper respository and version found.", imageRef: ImageReference{Repository: "turbonomic/zookeeper", Registries: []string{"registry.connect.redhat.com"}, Tag: "8.1.2", Sha: ""}, message: ""},
		{description: "Test cpopen/ibmcloud-object-storage-driver respository and sha found.", imageRef: ImageReference{Repository: "cpopen/ibmcloud-object-storage-driver", Registries: []string{"icr.io"}, Tag: "", Sha: "sha256:fc17bb3e89d00b3eb0f50b3ea83aa75c52e43d8e56cf2e0f17475e934eeeeb5f"}, message: ""},
		{description: "Test cpopen/ibmcloud-object-storage-plugin respository and sha found.", imageRef: ImageReference{Repository: "cpopen/ibmcloud-object-storage-plugin", Registries: []string{"icr.io"}, Tag: "", Sha: "sha256:cf654987c38d048bc9e654f3928e9ce9a2a4fd47ce0283bb5f339c1b99298e6e"}, message: ""},
	}

	for _, tc := range PassTestCases {
		t.Run(tc.description, func(t *testing.T) {
			found, err := IsImageInRegistry(tc.imageRef)
			require.NoError(t, err)
			require.True(t, found)
		})
	}

	FailTestCases := []testCase{
		{description: "Test postgresql-10-rhel7 version not found", imageRef: ImageReference{Repository: "rhscl/postgresql-10-rhel7", Registries: []string{"registry.access.redhat.com"}, Tag: "1.6.8", Sha: ""}, message: "Tag 1.6.8 not found"},
		{description: "Test rhel6.7 registry not found", imageRef: ImageReference{Repository: "rhel6.7", Registries: []string{"registry.notfound.com"}, Tag: "7.8", Sha: ""}, message: "No images found for Registry/Repository: registry.notfound.com/rhel6.7"},
		{description: "Test cpopen/ibmcloud-object-storage-plugin respository sha not found.", imageRef: ImageReference{Repository: "cpopen/ibmcloud-object-storage-plugin", Registries: []string{"icr.io"}, Tag: "", Sha: "sha256:ffff4987c38d048bc9e654f3928e9ce9a2a4fd47ce0283bb5f339c1b9929ffff"}, message: "Digest sha256:ffff4987c38d048bc9e654f3928e9ce9a2a4fd47ce0283bb5f339c1b9929ffff not found"},
	}

	for _, tc := range FailTestCases {
		t.Run(tc.description, func(t *testing.T) {
			found, err := IsImageInRegistry(tc.imageRef)
			require.Error(t, err)
			require.False(t, found)
			require.Contains(t, err.Error(), tc.message)
		})
	}
}
