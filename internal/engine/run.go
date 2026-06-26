package engine

import (
	"fmt"

	"github.com/allenbiji/clone-sage/internal/model"
	"github.com/allenbiji/clone-sage/internal/registry"
)

// ANSI Color Codes for pretty terminal output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

// Run executes the diagnostics. It returns true if the environment is healthy
// (no blockers failed), and false if a blocker failed.
func Run(cfg *model.ClonesageConfig) bool {
	fmt.Println(Cyan + "Running CloneSage Diagnostics...\n" + Reset)

	hasBlockerFailed := false
	passedCount := 0
	failedCount := 0

	// 1. Loop through every check in the YAML
	for _, checkCfg := range cfg.Checks {

		// 2. Ask the Registry to build the physical check
		check, err := registry.Build(checkCfg)
		if err != nil {
			// This shouldn't happen because of our Validation layer, but safe to catch!
			fmt.Printf("❌ %s [%s]: Internal Error - %v\n", checkCfg.Name, checkCfg.Type, err)
			hasBlockerFailed = true
			failedCount++
			continue
		}

		// 3. Execute the check!
		err = check.Execute()

		// 4. Handle the Result
		if err == nil {
			fmt.Printf("%s✅ %s%s\n", Green, checkCfg.Name, Reset)
			passedCount++
		} else {
			failedCount++

			// Evaluate Severity
			switch checkCfg.Severity {
			case model.SeverityInfo:
				fmt.Printf("%sℹ️  %s (Info)%s\n", Cyan, checkCfg.Name, Reset)
				fmt.Printf("   Reason: %v\n", err)
			case model.SeverityWarning:
				fmt.Printf("%s⚠️  %s (Warning)%s\n", Yellow, checkCfg.Name, Reset)
				fmt.Printf("   Reason: %v\n", err)
			case model.SeverityBlocker:
				fmt.Printf("%s❌ %s (BLOCKER)%s\n", Red, checkCfg.Name, Reset)
				fmt.Printf("   Reason: %v\n", err)
				if checkCfg.Message != "" {
					fmt.Printf("   Message: %s\n", checkCfg.Message)
				}
				if checkCfg.Fix != "" {
					fmt.Printf("   Fix: %s\n", checkCfg.Fix)
				}
				hasBlockerFailed = true
			}
		}
	}

	// 5. Print the Summary
	fmt.Println("\n----------------------------------------")
	if hasBlockerFailed {
		fmt.Printf("%s❌ DIAGNOSTICS FAILED: %d passed, %d failed%s\n", Red, passedCount, failedCount, Reset)
		return false
	}

	fmt.Printf("%s✅ DIAGNOSTICS PASSED: %d passed, %d failed (non-blocking)%s\n", Green, passedCount, failedCount, Reset)
	return true
}
