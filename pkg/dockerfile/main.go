package dockerfile

import (
	"bytes"
	_ "embed"
	"os"
	"strings"
	"text/template"

	"github.com/hammacktony/dlc/pkg/fileutils"
	"github.com/hammacktony/dlc/pkg/spec"
)

//go:embed Dockerfile.tmpl
var dockerfileTmpl string

// Create dockerfile from spec
func Create(spec *spec.Config) ([]byte, error) {
	t, err := template.New("Dockerfile.tmpl").Funcs(template.FuncMap{
		"getShortPythonVersion": func(version string) string {
			slice := strings.Split(version, ".")
			if len(slice) < 3 {
				return version
			}
			return strings.Join(slice[:2], ".")
		},
	}).Parse(dockerfileTmpl)
	if err != nil {
		return nil, err
	}

	var bytes bytes.Buffer
	err = t.Execute(&bytes, spec)
	if err != nil {
		return nil, err
	}

	return bytes.Bytes(), nil
}

// Write dockerfile to file
func WriteFile(path string, content []byte) error {
	// Output to stdout
	if path == "-" || path == "/dev/stdout" {
		return fileutils.WriteFile(os.Stdout, content)
	}
	// Output to file
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	return fileutils.WriteFile(f, content)
}
