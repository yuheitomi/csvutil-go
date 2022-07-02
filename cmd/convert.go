// Package cmd /*
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yuheitomi/csvutil/csvutil"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert CSV... schema",
	Short: "Convert CSV header as specified by schema",

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
		defer r.Close()

		schema, err := csvutil.ReadSchema(r)

		if err := convertFiles(csvFiles, schema, outputDirectory); err != nil {
			log.Fatal(err)
		}
	},
}

var outputDirectory = ""
var skipEmptyColumns = false

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&outputDirectory, "output", "o", "", "output directory")
	convertCmd.Flags().BoolVar(&skipEmptyColumns, "skip-empty", false, "Skip columns not renamed in schema")
}

func convertFiles(csvFiles []string, schema csvutil.HeaderSchema, outDir string) error {
	stat, err := os.Stat(outDir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("not a directory: %v", outDir)
	}

	for _, csvFile := range csvFiles {
		log.Printf("converting %v", csvFile)
		newFile := outputFilePath(csvFile, outDir)
		if err := convertCSV(csvFile, newFile, schema); err != nil {
			return err
		}
	}
	return nil
}

func convertCSV(csvFile string, outFile string, schema csvutil.HeaderSchema) error {
	r, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer r.Close()

	var ior io.Reader = r
	if isShiftJIS {
		ior = csvutil.ShiftJISEncoder(r)
	}

	w, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer w.Close()

	if err := csvutil.ConvertCSV(ior, w, schema, skipEmptyColumns); err != nil {
		return err
	}
	return nil
}

func outputFilePath(origFile string, outDir string) string {
	origName := filepath.Base(origFile)
	outFile := filepath.Join(outDir, origName)
	if filepath.Clean(origFile) == filepath.Clean(outFile) {
		log.Fatalf("same output file. %v -> %v", origFile, outFile)
	}
	return outFile
}
