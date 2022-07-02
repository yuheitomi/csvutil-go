package csvutil

import (
	"encoding/csv"
	"fmt"
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
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

func ConvertCSV(r io.Reader, w io.Writer, schema HeaderSchema, skipEmpty bool) error {
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

	newHeader := make([]string, 0)
	skipColumns := make([]bool, len(header))

	for i, col := range header {
		newCol, ok := schema[col]
		if ok {
			newHeader = append(newHeader, newCol)
		} else {
			if skipEmpty {
				skipColumns[i] = true
			} else {
				newHeader = append(newHeader, col)
			}
		}
	}
	if err := cw.Write(newHeader); err != nil {
		return err
	}

	// Just copy the following rows without changes.
	newRow := make([]string, 0, len(newHeader))
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if skipEmpty {
			for i, col := range row {
				if !skipColumns[i] {
					newRow = append(newRow, col)
				}
			}
			if err := cw.Write(newRow); err != nil {
				return err
			}
			newRow = newRow[:0]
		} else {
			if err := cw.Write(row); err != nil {
				return err
			}
		}
	}
	return nil
}

func ShiftJISEncoder(r io.Reader) *transform.Reader {
	return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
}
