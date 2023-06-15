package cli

import (
	"os"

	"github.com/hammacktony/dlc/pkg/dockerfile"
	"github.com/hammacktony/dlc/pkg/spec"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func BuildCmd() *cobra.Command {
	var exportFile string

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Use yaml to build a docker image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			containerSpec, err := spec.LoadConfig(args[0])
			if err != nil {
				zap.L().Error("failed to load container spec", zap.Error(err))
				os.Exit(1)
			}

			content, err := dockerfile.Create(containerSpec)
			if err != nil {
				zap.L().Error("failed to generate dockerfile", zap.Error(err))
				os.Exit(1)
			}

			if err := dockerfile.WriteFile(exportFile, content); err != nil {
				zap.L().Error("failed to export dockerfile", zap.Error(err))
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&exportFile, "export", "", "File to export dockerfile to")
	return cmd
}
