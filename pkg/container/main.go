package container

import "github.com/goccy/go-yaml"

type ContainerSpec struct {
	Cuda struct {
		Enabled bool   `yaml:"enabled"`
		Version string `yaml:"version"`
	} `yaml:"cuda"`
	SystemPackages []string `yaml:"system_packages"`
	Src            []string `yaml:"src"`
	Workdir        *string  `yaml:"workdir"`
	Python         struct {
		Version string `yaml:"version"`
		Poetry  *struct {
			Version string `yaml:"version"`
		} `yaml:"poetry"`
		ConfigSource []string `yaml:"config_source"`
	} `yaml:"python"`
}

func ReadSpec(fileBytes []byte) (*ContainerSpec, error) {
	var containerSpec ContainerSpec

	if err := yaml.Unmarshal(fileBytes, &containerSpec); err != nil {
		return nil, err
	}
	return &containerSpec, nil
}
