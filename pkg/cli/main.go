package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dlc",
	Short: "Making deep learning and containers easier to coexist",
	Long:  `Allow for easy construction of deep learning containers via a yaml file`,
}

func Execute() {
	rootCmd.AddCommand(BuildCmd())
	rootCmd.AddCommand(ShowCmd())
	rootCmd.AddCommand(GenerateCmd())
	rootCmd.AddCommand(ValidateCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
