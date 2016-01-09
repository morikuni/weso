package weso

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/peterh/liner"
)

type Template struct {
	templates map[string]*template.Template
}

func NewTemplateFile(path string) (*Template, error) {
	templates := make(map[string]*template.Template)
	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}
	for _, tt := range t.Templates() {
		if tt.Name() != path {
			templates[tt.Name()] = tt
		}
	}
	return &Template{templates}, nil
}

func (t *Template) IsDefined(name string) bool {
	_, ok := t.templates[name]
	return ok
}

func (t *Template) Names() []string {
	keys := make([]string, len(t.templates))
	for k := range t.templates {
		keys = append(keys, k)
	}
	return keys
}

func (t *Template) Apply(name string, args ...string) ([]byte, error) {
	vars := make(map[string]string)
	for i, a := range args {
		vars[fmt.Sprintf("_%d", i+1)] = a
	}
	buf := &bytes.Buffer{}
	if err := t.templates[name].Execute(buf, vars); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Completer(t *Template) liner.Completer {
	names := make([]string, len(t.Names()))
	for i, n := range t.Names() {
		names[i] = "." + n + " "
	}

	return func(line string) []string {
		var c []string
		for _, n := range names {
			if strings.HasPrefix(strings.ToLower(n), strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return c
	}
}
