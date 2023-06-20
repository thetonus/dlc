package spec

import (
	"github.com/goccy/go-yaml"
	"github.com/hammacktony/dlc/pkg/fileutils"
)

type SystemConfig struct {
	UbuntuVersion *string  `yaml:"ubuntu_version"`
	MambaVersion  *string  `yaml:"mamba_version"`
	Packages      []string `yaml:"packages"`
}

type CudaConfig struct {
	Enabled bool    `yaml:"enabled"`
	Version *string `yaml:"version"`
}

type PoetryConfig struct {
	Version string `yaml:"version"`
}

type PythonConfig struct {
	Version string        `yaml:"version"`
	Poetry  *PoetryConfig `yaml:"poetry"`
}

type ResourcesConfig struct {
	Config  []string `yaml:"config"`
	Test    []string `yaml:"test"`
	Project []string `yaml:"project"`
}

type Config struct {
	ProjectName string          `yaml:"project_name"`
	System      *SystemConfig   `yaml:"system"`
	Cuda        *CudaConfig     `yaml:"cuda"`
	Python      PythonConfig    `yaml:"python"`
	Resources   ResourcesConfig `yaml:"resources"`
}

// Given any type, return a reference to it.
func toPtr[T any](t T) *T {
	return &t
}

// Defaults
const (
	pythonVersion string = "3.10.11"
	ubuntuVersion string = "22.04"
	mambaVersion  string = "23.1.0-1"
	cudaEnabled   bool   = false
	cudaVersion   string = "11.8.0"
	poetryVersion string = "1.5.1"
)

func SetDefaults(config *Config) {
	if config.System == nil {
		config.System = &SystemConfig{}
	}

	if config.Cuda == nil {
		config.Cuda = &CudaConfig{Enabled: cudaEnabled, Version: nil}
	}

	// Default python package manager is poetry at the moment
	if config.Python.Poetry == nil {
		config.Python.Poetry = &PoetryConfig{Version: poetryVersion}
	}

	if config.System.UbuntuVersion == nil {
		config.System.UbuntuVersion = toPtr(ubuntuVersion)
	}

	if config.System.MambaVersion == nil {
		config.System.MambaVersion = toPtr(mambaVersion)
	}

	if config.Cuda.Enabled == true && config.Cuda.Version == nil {
		config.Cuda.Version = toPtr(cudaVersion)
	}
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (*Config, error) {
	var config Config

	bytes, err := fileutils.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	// Set default values
	SetDefaults(&config)

	return &config, nil
}

// GenerateConfig generates an example config file
func GenerateConfig(useCuda bool) *Config {
	config := Config{
		ProjectName: "example_project",
		Resources: ResourcesConfig{
			Config:  []string{"pyproject.toml", "poetry.lock"},
			Test:    []string{"test"},
			Project: []string{"example_project", "app"},
		},
		Python: PythonConfig{
			Version: "3.10.11",
		},
		System: &SystemConfig{
			Packages: []string{"ca-certificates", "tzdata"},
		},
	}

	if useCuda == true {
		config.Cuda = &CudaConfig{Enabled: true, Version: toPtr(cudaVersion)}
	}

	SetDefaults(&config)
	return &config
}

// Validates configuration
func ValidateConfig(path string) error {
	_, err := LoadConfig(path)
	return err
}
