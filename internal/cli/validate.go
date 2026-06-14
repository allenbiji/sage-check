package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewValidateCmd() *cobra.Command{
	validCmd := &cobra.Command{
		Use: "validate",
		Short: "Used to validate the checks in sage.yaml",
		Long: "This command can be to used to verify if the user-defined checks in sage.yaml are valid or not",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Validate command has been invoked")
			return nil
		},
	}

	return validCmd
}