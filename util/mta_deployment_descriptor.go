package util

type MtaDeploymentDescriptor struct {
	SchemaVersion string     `yaml:"_schema-version,omitempty"`
	ID            string     `yaml:"ID,omitempty"`
	Version       string     `yaml:"version,omitempty"`
	Modules       []Module   `yaml:"modules,omitempty"`
	Resources     []Resource `yaml:"resources,omitempty"`
}

type Module struct {
	Name                 string               `yaml:"name"`
	Type                 string               `yaml:"type"`
	Path                 string               `yaml:"path"`
	RequiredDependencies []RequiredDependency `yaml:"requires,omitempty"`
}

type Resource struct {
	Name       string                 `yaml:"name"`
	Type       string                 `yaml:"type"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

type RequiredDependency struct {
	Name       string                 `yaml:"name"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}
