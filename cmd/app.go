package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type HeaderSchema map[string]string

func generateTemplate(r io.Reader, w io.Writer) {
	cr := csv.NewReader(r)
	cw := csv.NewWriter(w)

	header, err := cr.Read()
	if err != nil {
		log.Fatal(err)
	}

	buf := []string{"", ""}
	for _, col := range header {
		buf[0] = col
		if err := cw.Write(buf); err != nil {
			log.Fatal(err)
		}
	}
	cw.Flush()
}

func readSchema(r io.Reader) (map[string]string, error) {
	cr := csv.NewReader(r)
	result := make(map[string]string)

	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(rec) >= 2 {
			if rec[1] != "" {
				result[rec[0]] = rec[1]
			}
		}
	}

	return result, nil
}

func convertCSVHeaderWithFiles(csvFiles []string, schema map[string]string, outDir string) error {
	stat, err := os.Stat(outDir)
	if err == os.ErrNotExist {
		return err
	}
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("not a directory: %v", outDir)
	}

	for _, csvFile := range csvFiles {
		log.Printf("converting %v", csvFile)
		r, err := os.Open(csvFile)
		if err != nil {
			return err
		}

		// TODO: use proper output file
		newFile := outputFilePath(csvFile, outDir)
		w, err := os.Create(newFile)
		if err != nil {
			return err
		}

		func(r *os.File, w *os.File) {
			defer r.Close()
			defer w.Close()
			if err := convertCSV(r, w, schema); err != nil {
				log.Fatal(err)
			}
		}(r, w)
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

func convertCSV(r io.Reader, w io.Writer, schema HeaderSchema) error {
	cr := csv.NewReader(r)
	cw := csv.NewWriter(w)

	header, err := cr.Read()
	if err != nil {
		return err
	}

	newHeader := make([]string, len(header))
	for i, col := range header {
		newCol, ok := schema[col]
		if ok {
			newHeader[i] = newCol
		} else {
			newHeader[i] = col
		}
	}
	cw.Write(newHeader)
	defer cw.Flush()

	// TODO: Reuse CSV buffer
	for {
		rows, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		cw.Write(rows)
	}

	return nil
}
