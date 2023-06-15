package spec

import (
	"github.com/goccy/go-yaml"
	"github.com/hammacktony/dlc/pkg/fileutils"
)

type SystemConfig struct {
	UbuntuVersion *string  `mapstructure:"ubuntu_version" yaml:"ubuntu_version"`
	MambaVersion  *string  `mapstructure:"mamba_version" yaml:"mamba_version"`
	Packages      []string `mapstructure:"packages" yaml:"packages"`
}

type CudaConfig struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
	Version string `mapstructure:"version" yaml:"version"`
}

type PoetryConfig struct {
	Version string `mapstructure:"version" yaml:"version"`
}

type PythonConfig struct {
	Version string        `mapstructure:"version" yaml:"version"`
	Poetry  *PoetryConfig `mapstructure:"poetry" yaml:"poetry"`
}

type ResourcesConfig struct {
	Config  []string `mapstructure:"config" yaml:"config"`
	Test    []string `mapstructure:"test" yaml:"test"`
	Project []string `mapstructure:"project" yaml:"project"`
}

type Config struct {
	ProjectName string          `mapstructure:"project_name" yaml:"project_name"`
	System      *SystemConfig   `mapstructure:"system" yaml:"system"`
	Cuda        *CudaConfig     `mapstructure:"cuda" yaml:"cuda"`
	Python      PythonConfig    `mapstructure:"python" yaml:"python"`
	Resources   ResourcesConfig `mapstructure:"resources" yaml:"resources"`
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
		config.Cuda = &CudaConfig{Enabled: cudaEnabled, Version: cudaVersion}
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
