package model

// this is the type for each of the checks
type CheckConfig struct {
	Name     string            `yaml:"name"`
	Type     string            `yaml:"type"`
	Severity string            `yaml:"severity"`
	Options  map[string]string `yaml:"options"`
	Message  string            `yaml:"message"`
	Why      string            `yaml:"why"`
	Fix      string            `yaml:"fix"`
}

//this is the type for the global clonesage config in its yaml file
type ClonesageConfig struct {
	Version  int                    `yaml:"version"`
	Defaults map[string]interface{} `yaml:"defaults"`
	Checks   []CheckConfig          `yaml:"checks"`
}
