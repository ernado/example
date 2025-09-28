package template

import (
	_ "embed"
	"io"
	"strings"
	"text/template"

	"github.com/ernado/example/internal/otel-iface-gen/internal/registry"
)

// Template is the instrumentation template. It is capable of generating the
// instrumentation wrapper for the given template.Data.
type Template struct {
	tmpl *template.Template
}

// New returns a new instance of Template.
func New() (Template, error) {
	tmpl, err := template.New("instrumentation").Funcs(templateFuncs).Parse(instrumentationTemplate)
	if err != nil {
		return Template{}, err
	}

	return Template{tmpl: tmpl}, nil
}

// Execute generates and writes the instrumentation implementation for the given
// data.
func (t Template) Execute(w io.Writer, data Data) error {
	return t.tmpl.Execute(w, data)
}

/*

Should generate following:

type FooInstrumentation struct {
	impl example.DB
	tracer trace.Tracer
	meter metric.Meter
}

func (i FooInstrumentation) CreateTask(ctx context.Context, title string) (ret *example.Task, err error) {
	return i.impl.CreateTask(ctx, title)
}

func (i FooInstrumentation) ListTasks(ctx context.Context) (ret []*example.Task, err error) {
	return i.impl.ListTasks(ctx)
}

func (i FooInstrumentation) DeleteTask(ctx context.Context, id int64) (err error) {
	return i.impl.DeleteTask(ctx, id)
}

func (i FooInstrumentation) GenerateError(ctx context.Context) (err error) {
	return i.impl.GenerateError(ctx)
}
*/

// instrumentationTemplate generates an instrumentation wrapper for an interface.
// language=GoTemplate
//
//go:embed template.tmpl
var instrumentationTemplate string

// This list comes from the golint codebase. Golint will complain about any of
// these being mixed-case, like "Id" instead of "ID".
var golintInitialisms = []string{
	"ACL", "API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS",
	"QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "TCP", "TLS", "TTL", "UDP", "UI", "UID", "UUID", "URI",
	"URL", "UTF8", "VM", "XML", "XMPP", "XSRF", "XSS",
}

var templateFuncs = template.FuncMap{
	"ImportStatement": func(imprt *registry.Package) string {
		if imprt.Alias == "" {
			return `"` + imprt.Path() + `"`
		}
		return imprt.Alias + ` "` + imprt.Path() + `"`
	},
	"SyncPkgQualifier": func(imports []*registry.Package) string {
		for _, imprt := range imports {
			if imprt.Path() == "sync" {
				return imprt.Qualifier()
			}
		}

		return "sync"
	},
	"Exported": func(s string) string {
		if s == "" {
			return ""
		}
		for _, initialism := range golintInitialisms {
			if strings.ToUpper(s) == initialism {
				return initialism
			}
		}
		return strings.ToUpper(s[0:1]) + s[1:]
	},
}
