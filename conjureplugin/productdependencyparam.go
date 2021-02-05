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
	"bytes"
	"fmt"
	"regexp"
	"text/template"
)

type RenderedProductDependencyProvider interface {
	RenderedProductDependencies() ([]RenderedProductDependency, error)
}

type renderedProductDependencyProviderImpl struct {
	params []ProductDependencyParam
	versionProvider VersionProvider
}

func (r *renderedProductDependencyProviderImpl) RenderedProductDependencies() ([]RenderedProductDependency, error) {
	var output []RenderedProductDependency
	for _, param := range r.params {
		rendered, err := param.RenderProductDependency(r.versionProvider)
		if err != nil {
			return nil, err
		}
		output = append(output, rendered)
	}
	return output, nil
}

// ProductDependencyParam represents a product dependency.
type ProductDependencyParam struct {
	ProductGroup string
	ProductName  string

	// MinimumVersion, MaximumVersion and RecommendedVersion are executed as Go templates to render the output version.
	// The rendered output of MinimumVersion and RecommendedVersion must be a valid orderable SLS version, while the
	// rendered output of MaximumVersion must be a valid SLS version matcher.
	//
	// The template environment supports the function "ProjectVersion", which renders the version of the enclosing
	// project. "ProjectVersion.(Major|Minor|Patch)" renders the major, minor or patch version of the product if the
	// project version matches the regular expression `^([0-9]+)\.([0-9]+)\.([0-9]+)` (if the project version does not
	// match this regular expression, executing these functions will result in an error).
	MinimumVersion     string
	MaximumVersion     string
	RecommendedVersion string
}

func (p *ProductDependencyParam) RenderProductDependency(versionProvider VersionProvider) (RenderedProductDependency, error) {
	renderedMinVersion, err := renderVersionTemplate(p.MinimumVersion, versionProvider)
	if err != nil {
		return RenderedProductDependency{}, err
	}
	if err := validateIsValidSLSVersion(renderedMinVersion); err != nil {
		return RenderedProductDependency{}, err
	}

	renderedMaxVersion, err := renderVersionTemplate(p.MaximumVersion, versionProvider)
	if err != nil {
		return RenderedProductDependency{}, err
	}
	if err := validateIsValidSLSMatcher(renderedMaxVersion); err != nil {
		return RenderedProductDependency{}, err
	}

	renderedRecommendedVersion, err := renderVersionTemplate(p.RecommendedVersion, versionProvider)
	if err != nil {
		return RenderedProductDependency{}, err
	}
	if renderedRecommendedVersion != "" {
		if err := validateIsValidSLSVersion(renderedRecommendedVersion); err != nil {
			return RenderedProductDependency{}, err
		}
	}

	return RenderedProductDependency{
		ProductGroup: p.ProductGroup,
		ProductName: p.ProductGroup,
		MinimumVersion: renderedMinVersion,
		MaximumVersion: renderedMaxVersion,
		RecommendedVersion: renderedRecommendedVersion,
	}, nil
}

// RenderedProductDependency represents a concrete product dependency.
type RenderedProductDependency struct {
	ProductGroup       string `json:"product-group" yaml:"product-group"`
	ProductName        string `json:"product-name" yaml:"product-name"`
	MinimumVersion     string `json:"minimum-version" yaml:"minimum-version"`
	MaximumVersion     string `json:"maximum-version" yaml:"maximum-version"`
	RecommendedVersion string `json:"recommended-version,omitempty" yaml:"recommended-version,omitempty"`
}

type templateProjectVersion string

func (v templateProjectVersion) String() string {
	return string(v)
}

func (v templateProjectVersion) Major() (string, error) {
	return v.partAtPos(0)
}

func (v templateProjectVersion) Minor() (string, error) {
	return v.partAtPos(1)
}

func (v templateProjectVersion) Patch() (string, error) {
	return v.partAtPos(2)
}

// orderableVersionRegexp is a regular expression that matches the major, minor and patch portions of an orderable
// version string. Intentionally omits a trailing "$" so that matcher will still match relevant portions of valid
// orderable versions like snapshots (for example, "1.0.0-1-gaaaaaaa").
var orderableVersionRegexp = regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9]+)`)

func (v templateProjectVersion) partAtPos(pos int) (string, error) {
	matches := orderableVersionRegexp.FindStringSubmatch(string(v))
	if matches == nil {
		return "", fmt.Errorf("version %q did not match regular expression for an orderable version", v)
	}
	if pos >= 3 {
		return "", fmt.Errorf("requested part at index %d, but valid orderable versions only have 3 parts", pos)
	}
	// pos+1 because element 0 contains entire match
	return matches[pos+1], nil
}

type VersionProvider interface {
	Version() (string, error)
}

// renderVersionTemplate renders the provided versionTemplateContent using the provided projectVersionProvider to
// retrieve the current version of the project. The template supports the function "ProjectVersion", which returns the
// result of executing the provided projectVersionProvider. The object returned by "ProjectVersion" renders the project
// version as a string and supports the functions "Major", "Minor" and "Patch" to return the major, minor and patch
// versions, respectively (unless the project version does not contain these components, in which case executing these
// functions returns an error).
func renderVersionTemplate(versionTemplateContent string, projectVersionProvider VersionProvider) (string, error) {
	tmpl := template.New("versionTemplate")
	tmpl.Funcs(template.FuncMap{
		"ProjectVersion": func() (templateProjectVersion, error) {
			version, err := projectVersionProvider.Version()
			if err != nil {
				return "", err
			}
			return templateProjectVersion(version), nil
		},
	})
	tmpl, err := tmpl.Parse(versionTemplateContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return "", err
	}
	return buf.String(), nil
}
