package csvutil

import (
	"encoding/csv"
	"fmt"
	"io"
)

type HeaderSchema map[string]string

func GenerateTemplate(r io.Reader, w io.Writer) error {
	cr := csv.NewReader(r)
	cw := csv.NewWriter(w)
	defer cw.Flush()

	header, err := cr.Read()
	if err != nil {
		return err
	}

	buf := []string{"", ""}
	for _, col := range header {
		buf[0] = col
		if err := cw.Write(buf); err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}
	}
	return nil
}

func ReadSchema(r io.Reader) (HeaderSchema, error) {
	cr := csv.NewReader(r)
	result := make(HeaderSchema)

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

func ConvertCSV(r io.Reader, w io.Writer, schema HeaderSchema) error {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	cw := csv.NewWriter(w)
	defer cw.Flush()

	// Read header and replace with new column names specified by schema.
	// Will use the original column name if replacement name is not provided.
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
	if err := cw.Write(newHeader); err != nil {
		return err
	}

	// Just copy the following rows without changes.
	for {
		rows, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if err := cw.Write(rows); err != nil {
			return err
		}
	}
	return nil
}
