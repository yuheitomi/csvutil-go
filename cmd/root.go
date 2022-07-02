package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csvutil",
	Short: "CSV file utility commands",
}

var isShiftJIS = false

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&isShiftJIS, "shiftjis", false, "Read as Shift-JIS file")
}
