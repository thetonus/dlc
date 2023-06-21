package cli

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/hammacktony/dlc/pkg/fileutils"
	"github.com/hammacktony/dlc/pkg/spec"
	"github.com/spf13/cobra"
)

func GenerateCmd() *cobra.Command {
	var useCuda bool

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate example config",
		Run: func(cmd *cobra.Command, args []string) {
			config := spec.GenerateConfig(useCuda)
			content, err := yaml.Marshal(config)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if err := fileutils.WriteFile(os.Stdout, content); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&useCuda, "cuda", false, "generate cuda-enabled config")
	return cmd
}
