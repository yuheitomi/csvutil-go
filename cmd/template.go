// Package cmd /*
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"yuheitomi/csvutil/csvutil"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Generate template file(s) for specifying new column names and types from CSV(s)",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("needs to specify a CSV")
			return
		}

		filename := args[0]
		if ext := filepath.Ext(filename); ext != ".csv" {
			fmt.Printf("Not a CSV file: %v", filename)
			return
		}

		r, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		if err := csvutil.GenerateTemplate(r, os.Stdout); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
