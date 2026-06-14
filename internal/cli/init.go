package cli

import (
	"fmt"
	"os"

	"github.com/allenbiji/clone-sage/internal/detect"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Run init to initalize the project and generate sage-auto.yml",
		Long:  "Run this command to scan your entire repository and from the inferences in your repo, a sage-auto.yml file will be generated and which can also be extended via a sage.yml file",
		RunE: func(cmd *cobra.Command, args []string) error {

			// Run ScanRepo() to generate baseline 
			cfgs := detect.ScanRepo()

			if len(cfgs.Checks) == 0 {
				fmt.Println("No recognised frameworks found. Generating empty baseline.")
			} else {
				fmt.Printf("Detected %d requirements. Building configuration...\n", len(cfgs.Checks))
			}

			// Marshal the struct cleanly into YAML bytes
			configYaml, err := yaml.Marshal(cfgs)
			if err != nil {
				return err
			}

			// Write generated baseline into sage-auto.yml file with standard 0644 file permissions
			err = os.WriteFile("sage-auto.yml", configYaml, 0644)
			if err != nil {
				return err
			} 

			fmt.Println("Baseline has been generated successfully in sage-auto.yml!")
			fmt.Println("Run sage check to verify your local environment")
			return nil
		},
	}

	return initCmd
}
