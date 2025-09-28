// Binary otel-sort sorts OpenTelemetry registry bundles by group ID.
//
// https://github.com/open-telemetry/weaver/issues/964
package main

import (
	"bytes"
	"flag"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

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

	var r map[string][]map[string]any
	if err := yaml.Unmarshal(data, &r); err != nil {
		panic(err)
	}

	// Sort groups by ID.
	slices.SortFunc(r["groups"], func(a, b map[string]any) int {
		return strings.Compare(a["id"].(string), b["id"].(string))
	})

	out := new(bytes.Buffer)
	e := yaml.NewEncoder(out)
	e.SetIndent(2)
	if err := e.Encode(r); err != nil {
		panic(err)
	}

	//#nosec: G306
	if err := os.WriteFile(arg.File, out.Bytes(), 0o644); err != nil {
		panic(err)
	}
}
