package main

import (
	"errors"
	"os"
	"strings"

	"github.com/json-iterator/go"
	"github.com/yujinlin0224/twse-isin-tables-crawler"
)

func main() {
	var err error
	tables, errs := twseisintablescrawler.CrawlAll()
	if len(errs) > 0 {
		errStrings := make([]string, len(errs))
		for i := range errs {
			errStrings[i] = errs[i].Error()
		}
		panic(errors.New(strings.Join(errStrings, "\n")))
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
