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
	"fmt"
	"regexp"
)

var (
	regexpSLSVersionRelease                  = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	regexpSLSVersionReleaseSnapshot          = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+-[0-9]+-g[a-f0-9]+$`)
	regexpSLSVersionReleaseCandidate         = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+-rc[0-9]+$`)
	regexpSLSVersionReleaseCandidateSnapshot = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+-rc[0-9]+-[0-9]+-g[a-f0-9]+$`)
	regexpSLSVersionMatcher                  = regexp.MustCompile(`^((x\.x\.x)|([0-9]+\.x\.x)|([0-9]+\.[0-9]+\.x)|([0-9]+\.[0-9]+\.[0-9]+))$`)
)

func validateIsValidSLSVersion(in string) error {
	if isValidSLSVersion(in) {
		return nil
	}
	return fmt.Errorf("%q is not a valid SLS version", in)
}

func isValidSLSVersion(in string) bool {
	return regexpSLSVersionRelease.MatchString(in) ||
		regexpSLSVersionReleaseSnapshot.MatchString(in) ||
		regexpSLSVersionReleaseCandidate.MatchString(in) ||
		regexpSLSVersionReleaseCandidateSnapshot.MatchString(in)
}

func validateIsValidSLSMatcher(in string) error {
	if isValidSLSVersionMatcher(in) {
		return nil
	}
	return fmt.Errorf("%q is not a valid SLS version matcher", in)
}

func isValidSLSVersionMatcher(in string) bool {
	return regexpSLSVersionMatcher.MatchString(in)
}
