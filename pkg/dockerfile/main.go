package dockerfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hammacktony/dlc/pkg/container"
)

// TODO: Use runes for tabs and new lines (maybe for Dockerfile keywords)

const (
	ubuntuVersion   = "22.04"
	ubuntuBaseImage = "ubuntu:%s"
	nvidiaBaseImage = "nvidia/cuda:%s-runtime-ubuntu%s"
)

func baseImage(dockerFile *strings.Builder, spec *container.ContainerSpec) {
	if spec.Cuda.Enabled == true {
		dockerFile.WriteString(fmt.Sprintf("FROM "+nvidiaBaseImage+" AS base\n", spec.Cuda.Version, ubuntuVersion))
	} else {
		dockerFile.WriteString(fmt.Sprintf("FROM "+ubuntuBaseImage+" AS base\n", ubuntuVersion))
	}
}

func installSystemPackages(dockerFile *strings.Builder, spec *container.ContainerSpec) {
	var installerSteps strings.Builder

	installerSteps.WriteString("RUN apt-get update \\" + "\n")
	installerSteps.WriteString("\t&& apt-install -y --no-install-recommends \\\n")

	for _, pkg := range spec.SystemPackages {
		installerSteps.WriteString("\t" + pkg + " \\\n")
	}

	installerSteps.WriteString("\t&& rm -rf /var/lib/apt/lists/*\n")
	dockerFile.WriteString(installerSteps.String())
}

// Create dockerfile from spec
func Create(spec *container.ContainerSpec) string {
	var dockerFile strings.Builder

	baseImage(&dockerFile, spec)
	installSystemPackages(&dockerFile, spec)

	if spec.Workdir != nil {
		fmt.Printf("WORKDIR %s\n", *spec.Workdir)
	}
	for _, item := range spec.Python.ConfigSource {
		dockerFile.WriteString("COPY " + item + " " + item + "\n")
	}

	return dockerFile.String()
}

func writeFile(f io.Writer, content string) error {
	writer := bufio.NewWriter(f)
	_, err := writer.WriteString(content)
	defer writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

// Write dockerfile to file
func WriteFile(path string, content string) error {
	// Output to stdout
	if path == "-" || path == "/dev/stdout" {
		return writeFile(os.Stdout, content)
	}
	// Output to file
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	return writeFile(f, content)
}
