package cli

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/hammacktony/dlc/pkg/fileutils"
	"github.com/hammacktony/dlc/pkg/spec"
	"github.com/spf13/cobra"
)

func ShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show full config spec",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := spec.LoadConfig(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

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
	return cmd
}
