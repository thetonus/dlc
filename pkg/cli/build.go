package cli

import (
	"fmt"
	"os"

	"github.com/hammacktony/dlc/pkg/container"
	"github.com/hammacktony/dlc/pkg/dockerfile"
	"github.com/hammacktony/dlc/pkg/fileutils"
	"github.com/spf13/cobra"
)

func BuildCmd() *cobra.Command {
	var exportFile string

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Use yaml to build a docker image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file, err := fileutils.ReadFile(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			containerSpec, err := container.ReadSpec(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			content := dockerfile.Create(containerSpec)
			if err := dockerfile.WriteFile(exportFile, content); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&exportFile, "export", "", "File to export dockerfile to")
	return cmd
}
