package main

import (
	"errors"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func main() {
	var err error
	tables := make([]Table, 0)
	for i := 1; i <= 12; i++ {
		t, errs := crawl(i)
		if len(errs) > 0 {
			errStrings := make([]string, len(errs))
			for i := range errs {
				errStrings[i] = errs[i].Error()
			}
			panic(errors.New(strings.Join(errStrings, "\n")))
		}
		tables = append(tables, t...)
	}
	jsonFile, err := os.Create("tables.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	err = jsoniter.NewEncoder(jsonFile).Encode(&tables)
	if err != nil {
		panic(err)
	}
}
