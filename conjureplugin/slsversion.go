package conjureplugin

import (
	"fmt"
	"regexp"
)

var (
	regexpSLSVersionRelease = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	regexpSLSVersionReleaseSnapshot = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+-[0-9]+-g[a-f0-9]+$`)
	regexpSLSVersionReleaseCandidate = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+-rc[0-9]+$`)
	regexpSLSVersionReleaseCandidateSnapshot = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+-rc[0-9]+-[0-9]+-g[a-f0-9]+$`)
	regexpSLSVersionMatcher = regexp.MustCompile(`^((x\.x\.x)|([0-9]+\.x\.x)|([0-9]+\.[0-9]+\.x)|([0-9]+\.[0-9]+\.[0-9]+))$`)
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