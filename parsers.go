package twseisintablescrawler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/json-iterator/go"
	"golang.org/x/exp/slices"
)

type Parser func(Language, string) interface{}

var (
	StringParser = func(_ Language, s string) interface{} {
		return strings.TrimSpace(s)
	}
	NumberParser = func(_ Language, s string) interface{} {
		number, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
		return number
	}
	DateParser = func(_ Language, s string) interface{} {
		dateString := strings.ReplaceAll(strings.TrimSpace(s), "/", "-")
		if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(dateString) {
			return nil
		}
		date, _ := civil.ParseDate(dateString)
		return date
	}
	MultiLanguageTextParser = func(language Language, s string) interface{} {
		return NewMultiLanguageText(language, strings.TrimSpace(s))
	}
)

func ParseColumnsAndRow(language Language, columnLabels []string, tableRow []string) ([]Column, *Row, error) {
	if tableRow != nil && len(columnLabels) != len(tableRow) {
		return nil, nil, fmt.Errorf("table row length not match: %d != %d", len(tableRow), len(columnLabels))
	}
	rawRow := make(map[string]interface{})
	columns := make([]Column, 0, len(columnLabels)+2)
	for i := range columnLabels {
		columnLabel := strings.TrimSpace(columnLabels[i])
		var tableData string
		if tableRow != nil {
			tableData = tableRow[i]
		}
		supportedColumnIndex := slices.IndexFunc(SupportedColumns, func(column Column) bool {
			return column.Label.Chinese == columnLabel || column.Label.English == columnLabel
		})
		if supportedColumnIndex == -1 {
			return nil, nil, fmt.Errorf("column %s not supported", columnLabel)
		}
		column := SupportedColumns[supportedColumnIndex]
		if strings.Contains(column.Key, ",") {
			splitColumnKey := strings.Split(column.Key, ",")
			if len(splitColumnKey) != 2 {
				return nil, nil, fmt.Errorf("column %s key format error: %s", columnLabel, column.Key)
			}
			var splitColumnLabel []string
			if strings.Contains(columnLabel, "&") {
				splitColumnLabel = strings.Split(columnLabel, "&")
				if len(splitColumnLabel) != 2 {
					return nil, nil, fmt.Errorf("column %s label format error: %s", columnLabel, column.Label)
				}
			} else if strings.Contains(columnLabel, "及") {
				splitColumnLabel = strings.Split(columnLabel, "及")
				if len(splitColumnLabel) != 2 || len(splitColumnLabel[0]) <= len(splitColumnLabel[1]) {
					return nil, nil, fmt.Errorf("column %s label format error: %s", columnLabel, column.Label)
				}
				splitColumnLabel[1] = splitColumnLabel[0][:len(splitColumnLabel[0])-len(splitColumnLabel[1])] + splitColumnLabel[1]
			}
			var splitTableData []string
			if tableRow != nil {
				splitTableData = strings.SplitN(tableData, "　", len(splitColumnKey))
				if len(splitTableData) != 2 {
					return nil, nil, fmt.Errorf("column %s value format error: %s", columnLabel, tableData)
				}
			}
			for j := range splitColumnKey {
				if tableRow != nil {
					rawRow[splitColumnKey[j]] = column.parsers[j](language, splitTableData[j])
				}
				columnLabel := strings.TrimSpace(splitColumnLabel[j])
				columns = append(columns, Column{
					Key:   splitColumnKey[j],
					Label: NewMultiLanguageText(language, columnLabel),
				})
			}
		} else {
			if tableRow != nil {
				rawRow[column.Key] = column.parsers[0](language, tableData)
			}
			columns = append(columns, column)
		}
	}
	var row *Row
	if tableRow != nil {
		rowJSON, err := jsoniter.Marshal(rawRow)
		if err != nil {
			return nil, nil, err
		}
		row = new(Row)
		err = jsoniter.Unmarshal(rowJSON, row)
		if err != nil {
			return nil, nil, err
		}
	}
	return columns, row, nil
}
