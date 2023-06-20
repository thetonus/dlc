package cli

import (
	"os"

	"github.com/hammacktony/dlc/pkg/spec"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func ValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate config spec",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := spec.ValidateConfig(args[0])
			if err != nil {
				zap.L().Error("validation failed", zap.Error(err))
				os.Exit(1)
			}
			zap.L().Info("validation successful")
		},
	}
	return cmd
}
