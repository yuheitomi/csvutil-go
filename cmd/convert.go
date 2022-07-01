/*
Copyright © 2022 Yuhei Kuratomi

*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [CSV files] [schema]",
	Short: "Convert a CSV header referring to the specified schema",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatalf("Needs at least one CSV file and Schema file")
		}

		size := len(args)
		csvFiles := args[0 : size-1]
		schemaFile := args[size-1]

		r, err := os.Open(schemaFile)
		if err != nil {
			log.Fatal(err)
		}

		schema, err := readSchema(r)

		if err := convertCSVHeaderWithFiles(csvFiles, schema, outputDirectory); err != nil {
			log.Fatal(err)
		}
	},
}

var outputDirectory = ""

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&outputDirectory, "output", "o", "", "output directory")
}
