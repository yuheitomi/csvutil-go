package csvutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type HeaderSchema map[string]Schema

type Schema struct {
	Name string
	Type string
}

var dateFormatRegexP *regexp.Regexp = regexp.MustCompile("DATE\\(([A-z-/:]+)\\)")

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
	cr.FieldsPerRecord = -1 // no check to capture optional third column
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
				f1 := rec[1]
				f2 := ""
				if len(rec) >= 3 {
					f2 = rec[2]
				}
				result[rec[0]] = Schema{
					Name: f1,
					Type: f2,
				}
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
	headerTypes := make([]string, len(header))
	skipColumns := make([]bool, len(header))

	for i, col := range header {
		newCol, ok := schema[col]
		if ok {
			newHeader = append(newHeader, newCol.Name)
			headerTypes[i] = newCol.Type
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

		for i, col := range row {
			if !skipColumns[i] {
				if strings.Contains(headerTypes[i], "DATE") {
					format := getDateFormat(headerTypes[i])
					col = convertDate(col, format)
				}
				newRow = append(newRow, col)
			}
		}
		if err := cw.Write(newRow); err != nil {
			return err
		}
		newRow = newRow[:0]
	}
	return nil
}

func getDateFormat(fmtString string) string {
	subMatches := dateFormatRegexP.FindSubmatch([]byte(fmtString))
	if len(subMatches) < 1 {
		return "2006-01-02"
	}

	format := string(subMatches[1])
	format = strings.Replace(format, "YYYY", "2006", 1)
	format = strings.Replace(format, "MM", "01", 1)
	format = strings.Replace(format, "DD", "02", 1)
	return string(format)
}

func convertDate(field string, dateFmt string) string {
	date, err := time.ParseInLocation(dateFmt, field, time.Local)
	if err != nil {
		log.Println(err)
		return ""
	}
	return date.String()
}

func ShiftJISEncoder(r io.Reader) *transform.Reader {
	return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
}
