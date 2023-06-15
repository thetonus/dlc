package cli

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

var rootCmd = &cobra.Command{
	Use:   "dlc",
	Short: "Making deep learning and containers easier to coexist",
	Long:  `Allow for easy construction of deep learning containers via a yaml file`,
}

func Execute() {
	rootCmd.AddCommand(BuildCmd())
	rootCmd.AddCommand(ShowCmd())
	rootCmd.AddCommand(GenerateCommand())
	if err := rootCmd.Execute(); err != nil {
		zap.L().Error("Error executing command", zap.Error(err))
		os.Exit(1)
	}
}
