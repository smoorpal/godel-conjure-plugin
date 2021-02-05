// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package conjureplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type versionProviderFn func() (string, error)

func (fn versionProviderFn) Version() (string, error) {
	return fn()
}

func TestRenderVersionTemplate(t *testing.T) {
	for _, tc := range []struct {
		name         string
		version      string
		testTemplate string
		want         string
		wantError    string
	}{
		{
			name:         "renders string literal",
			version:      "1.2.3",
			testTemplate: "2.3.4",
			want:         "2.3.4",
		},
		{
			name:         "render projectVersion in template",
			version:      "1.2.3",
			testTemplate: "{{ProjectVersion}}",
			want:         "1.2.3",
		},
		{
			name:         "render projectVersion with snapshot version",
			version:      "1.0.0-1-gaaaaaaa",
			testTemplate: "{{ProjectVersion}}",
			want:         "1.0.0-1-gaaaaaaa",
		},
		{
			name:         "Major function returns major version",
			version:      "1.2.3",
			testTemplate: "{{ProjectVersion.Major}}.x.x",
			want:         "1.x.x",
		},
		{
			name:         "Minor function returns minor version",
			version:      "1.2.3",
			testTemplate: "1.{{ProjectVersion.Minor}}.x",
			want:         "1.2.x",
		},
		{
			name:         "Patch function returns patch version",
			version:      "1.2.3",
			testTemplate: "1.2.{{ProjectVersion.Patch}}",
			want:         "1.2.3",
		},
		{
			name:         "Major, minor and patch functions work with RC snapshot version",
			version:      "1.2.3-rc1-1-gaaaaaaa",
			testTemplate: "{{ProjectVersion.Major}}.{{ProjectVersion.Minor}}.{{ProjectVersion.Patch}}",
			want:         "1.2.3",
		},
		{
			name:         "Major function errors if version does not match orderable regexp",
			version:      "unspecified",
			testTemplate: "{{ProjectVersion.Major}}.x.x",
			wantError:    `template: versionTemplate:1:16: executing "versionTemplate" at <ProjectVersion.Major>: error calling Major: version "unspecified" did not match regular expression for an orderable version`,
		},
		{
			name:         "Minor function errors if version does not match orderable regexp",
			version:      "unspecified",
			testTemplate: "1.{{ProjectVersion.Minor}}.x",
			wantError:    `template: versionTemplate:1:18: executing "versionTemplate" at <ProjectVersion.Minor>: error calling Minor: version "unspecified" did not match regular expression for an orderable version`,
		},
		{
			name:         "Patch function errors if version does not match orderable regexp",
			version:      "unspecified",
			testTemplate: "1.2.{{ProjectVersion.Patch}}",
			wantError:    `template: versionTemplate:1:20: executing "versionTemplate" at <ProjectVersion.Patch>: error calling Patch: version "unspecified" did not match regular expression for an orderable version`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			versionProvider := versionProviderFn(func() (string, error) {
				return tc.version, nil
			})
			got, err := renderVersionTemplate(tc.testTemplate, versionProvider)
			if tc.wantError != "" {
				assert.EqualError(t, err, tc.wantError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
