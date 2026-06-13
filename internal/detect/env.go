package detect

import (
	"strings"

	"github.com/allenbiji/clone-sage/internal/model"
)

func detectEnv() []model.CheckConfig{
	var checks []model.CheckConfig

	if fileExists(".env.example"){
		checks = append(checks, model.CheckConfig{
			Name: "env-file-exists",
			Type: "file_exists",
			Severity: "blocker",
			Options: map[string]string{"path":".env"},
			Message: "You must create a .env file",
			Fix: "Run cp .env.example .env",
		})

		keys := extractEnvKeys(".env.example")

		for _, key := range keys{
			checks = append(checks, model.CheckConfig{
				Name: strings.ToLower(key) + "-configured",
				Type: "env_exists",
				Severity: "blocker",
				Options: map[string]string{"key": key},
				Message: key + "is missing from your envirnonment",
			})
		}
	}

	return checks
}