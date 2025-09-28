// Binary otel-sort sorts OpenTelemetry registry bundles by group ID.
//
// https://github.com/open-telemetry/weaver/issues/964
package main

import (
	"flag"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type Registry struct {
	Groups []Group `yaml:"groups"`
}

type Group struct {
	ID         string      `yaml:"id"`
	Type       string      `yaml:"type"`
	Brief      string      `yaml:"brief"`
	Stability  string      `yaml:"stability,omitempty"`
	Attributes []Attribute `yaml:"attributes,omitempty"`
	MetricName string      `yaml:"metric_name,omitempty"`
	Instrument string      `yaml:"instrument,omitempty"`
	Unit       string      `yaml:"unit,omitempty"`
	Lineage    *Lineage    `yaml:"lineage,omitempty"`
}

type Attribute struct {
	Name             string   `yaml:"name"`
	Type             string   `yaml:"type"`
	Brief            string   `yaml:"brief"`
	Examples         []string `yaml:"examples,omitempty"`
	RequirementLevel string   `yaml:"requirement_level"`
	Stability        string   `yaml:"stability"`
}

type Lineage struct {
	Provenance Provenance             `yaml:"provenance"`
	Attributes map[string]LineageAttr `yaml:"attributes,omitempty"`
}

type Provenance struct {
	RegistryID string `yaml:"registry_id"`
	Path       string `yaml:"path"`
}

type LineageAttr struct {
	SourceGroup             string   `yaml:"source_group"`
	InheritedFields         []string `yaml:"inherited_fields"`
	LocallyOverriddenFields []string `yaml:"locally_overridden_fields"`
}

func main() {
	var arg struct {
		File string
	}
	flag.StringVar(&arg.File, "f", "", "Path to the registry file.")
	flag.Parse()

	data, err := os.ReadFile(arg.File)
	if err != nil {
		panic(err)
	}

	var r Registry
	if err := yaml.Unmarshal(data, &r); err != nil {
		panic(err)
	}

	// Sort groups by ID.
	slices.SortFunc(r.Groups, func(a, b Group) int {
		return strings.Compare(a.ID, b.ID)
	})

	out, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(arg.File, out, 0o644); err != nil {
		panic(err)
	}
}
