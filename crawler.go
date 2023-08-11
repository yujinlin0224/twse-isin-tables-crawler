package twseisintablescrawler

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/gocolly/colly"
	"github.com/jinzhu/copier"
)

const TableDomain = "isin.twse.com.tw"

var TableURLFormatOfLanguages = map[Language]string{
	LanguageEnglish: `https://isin.twse.com.tw/isin/e_C_public.jsp?strMode=%d`,
	LanguageChinese: `https://isin.twse.com.tw/isin/C_public.jsp?strMode=%d`,
}

var TableIndexes = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

func Crawl(tableIndex int) ([]Table, []error) {
	var errs []error

	if !slices.Contains(TableIndexes, tableIndex) {
		return nil, []error{fmt.Errorf("table index %d is not supported", tableIndex)}
	}

	c := colly.NewCollector(
		colly.AllowedDomains(TableDomain),
		colly.Async(true),
		colly.MaxDepth(1),
	)

	c.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "big5"
		fmt.Println("Visiting", r.URL, "at", time.Now())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL, "at", time.Now())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL, "at", time.Now())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error", r.Request.URL, err)
		errs = append(errs, err)
	})

	tableTemplateOfLanguages := map[Language]*Table{}
	tablesOfLanguages := map[Language]*[]Table{}
	for _, language := range []Language{LanguageEnglish, LanguageChinese} {
		tableTemplateOfLanguages[language] = new(Table)
		tableTemplateOfLanguages[language].Index = tableIndex
		tables := make([]Table, 0)
		tablesOfLanguages[language] = &tables
	}

	crawlTable := func(language Language, tableTemplate *Table, tables *[]Table) {
		tableURL := fmt.Sprintf(TableURLFormatOfLanguages[language], tableIndex)
		c.OnHTML("h2 > strong > font.h1:not(:has(center))", func(e *colly.HTMLElement) {
			if e.Request.URL.String() != tableURL {
				return
			}
			tableTemplate.Title = NewMultiLanguageText(language, e.Text)
		})
		c.OnHTML("h2 > strong > font.h1 > center", func(e *colly.HTMLElement) {
			if e.Request.URL.String() != tableURL {
				return
			}
			rawDate := strings.Split(e.Text, ":")[1]
			tableTemplate.UpdatedDate = DateParser(language, rawDate).(civil.Date)
		})
		c.OnHTML("table.h4 > tbody", func(e *colly.HTMLElement) {
			if e.Request.URL.String() != tableURL {
				return
			}
			rows := make([][]string, 0)
			tableRows := e.DOM.Find("tr").Nodes
			for _, tableRow := range tableRows {
				row := make([]string, 0)
				for tableData := tableRow.FirstChild; tableData != nil; tableData = tableData.NextSibling {
					if tableData.Data != "td" {
						continue
					}
					if tableData.FirstChild != nil {
						if tableData.FirstChild.Data != "b" {
							row = append(row, tableData.FirstChild.Data)
						} else {
							row = append(row, tableData.FirstChild.FirstChild.Data)
						}
					} else {
						row = append(row, "")
					}
				}
				rows = append(rows, row)
			}
			var columnLabels []string
			var subtitle *MultiLanguageText
			var table *Table
			for i, row := range rows {
				if i == 0 {
					columnLabels = row
					continue
				}
				if len(row) == 1 {
					if table != nil {
						*tables = append(*tables, *table)
						table = nil
					}
					continue
				}
				if table == nil {
					table = new(Table)
					table.Subtitle = subtitle
				}
				columns, row, err := ParseColumnsAndRow(language, columnLabels, row)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				if table.Columns == nil {
					table.Columns = columns
				}
				if table.Rows == nil {
					table.Rows = make([]Row, 0)
				}
				table.Rows = append(table.Rows, *row)
			}
			if table != nil {
				*tables = append(*tables, *table)
			}
			if len(*tables) == 0 {
				columns, _, err := ParseColumnsAndRow(language, columnLabels, nil)
				if err != nil {
					errs = append(errs, err)
					return
				}
				table = new(Table)
				table.Columns = columns
				table.Rows = make([]Row, 0)
				*tables = append(*tables, *table)
			}
		})

		c.Visit(tableURL)
	}

	for _, language := range languages {
		crawlTable(language, tableTemplateOfLanguages[language], tablesOfLanguages[language])
	}
	c.Wait()
	for i := 1; i < len(languages); i++ {
		if len(*tablesOfLanguages[languages[i]]) != len(*tablesOfLanguages[languages[i-1]]) {
			errs = append(errs, fmt.Errorf("table count mismatch between %s and %s", languages[i], languages[i-1]))
		}
	}
	for _, language := range languages {
		tableTemplate := tableTemplateOfLanguages[language]
		for i, table := range *tablesOfLanguages[language] {
			copier.CopyWithOption(&table, &tableTemplate, copier.Option{IgnoreEmpty: true, DeepCopy: true})
			(*tablesOfLanguages[language])[i] = table
		}
	}
	tables := make([]Table, len(*tablesOfLanguages[languages[0]]))
	for _, language := range languages {
		for i, tableOfLanguage := range *tablesOfLanguages[language] {
			table := tables[i]
			copier.CopyWithOption(&table, &tableOfLanguage, copier.Option{IgnoreEmpty: true, DeepCopy: true})
			tables[i] = table
		}
	}
	return tables, errs
}

func CrawlAll() ([]Table, []error) {
	var tables []Table
	var errs []error
	for _, tableIndex := range TableIndexes {
		t, e := Crawl(tableIndex)
		tables = append(tables, t...)
		errs = append(errs, e...)
	}
	return tables, errs
}
