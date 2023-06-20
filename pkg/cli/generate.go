package cli

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/hammacktony/dlc/pkg/fileutils"
	"github.com/hammacktony/dlc/pkg/spec"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
				zap.L().Error("failed to marshal config", zap.Error(err))
				os.Exit(1)
			}

			if err := fileutils.WriteFile(os.Stdout, content); err != nil {
				zap.L().Error("failed to write config", zap.Error(err))
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&useCuda, "cuda", false, "generate cuda-enabled config")
	return cmd
}
