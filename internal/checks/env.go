package checks

import (
	"fmt"

	"github.com/allenbiji/clone-sage/internal/detect"
	"github.com/allenbiji/clone-sage/internal/model"
	"github.com/allenbiji/clone-sage/internal/registry"
)

type EnvCheck struct {
	Key    string
	EnvMap map[string]string
}

func (e *EnvCheck) Execute() error {
	val, exists := e.EnvMap[e.Key]
	if !exists {
		return fmt.Errorf("key %q not found in .env", e.Key)
	}
	if val == "" {
		return fmt.Errorf("key %q is in .env but has no value", e.Key)
	}
	return nil
}

// cachedEnvMap holds the parsed .env contents for the lifetime of the process.
// The factory populates it on first call; all subsequent env_exists checks reuse it.
var cachedEnvMap map[string]string

func buildEnvExistsCheck(cfg model.CheckConfig) (registry.Check, error) {
	key, ok := cfg.Options["key"]
	if !ok || key == "" {
		return nil, fmt.Errorf("env_exists check requires a 'key' option")
	}

	if cachedEnvMap == nil {
		m, err := detect.ExtractEnvKeys(".env")
		if err != nil {
			return nil, fmt.Errorf("could not read .env: %w", err)
		}
		cachedEnvMap = m
	}

	return &EnvCheck{Key: key, EnvMap: cachedEnvMap}, nil
}

func init() {
	registry.Register(model.TypeEnvExists, buildEnvExistsCheck)
}
