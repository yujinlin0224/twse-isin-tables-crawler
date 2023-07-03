package main

import (
	"cloud.google.com/go/civil"
)

type Language string

const (
	LanguageEnglish Language = "english"
	LanguageChinese Language = "chinese"
)

var languages = []Language{
	LanguageEnglish,
	LanguageChinese,
}

type MultiLanguageText struct {
	English string `json:"english"`
	Chinese string `json:"chinese"`
}

func NewMultiLanguageText(language Language, text string) MultiLanguageText {
	switch language {
	case LanguageEnglish:
		return MultiLanguageText{English: text}
	case LanguageChinese:
		return MultiLanguageText{Chinese: text}
	}
	return MultiLanguageText{}
}

type Column struct {
	Key     string            `json:"key"`
	Label   MultiLanguageText `json:"label"`
	parsers []Parser
}

func NewColumn(key, englishLabel, chineseLabel string, parsers ...Parser) Column {
	return Column{
		Key: key,
		Label: MultiLanguageText{
			English: englishLabel,
			Chinese: chineseLabel,
		},
		parsers: parsers,
	}
}

var supportedColumns = []Column{
	NewColumn("securityCode,securityName", "Security Code & Security Name", "有價證券代號及名稱", stringParser, multiLanguageTextParser),
	NewColumn("isinCode", "ISIN Code", "國際證券辨識號碼(ISIN Code)", stringParser),
	NewColumn("issuedDate", "Date Issued", "公開發行日", dateParser),
	NewColumn("industrialGroup", "Industrial Group", "產業別", multiLanguageTextParser),
	NewColumn("cfiCode", "CFICode", "CFICode", stringParser),
	NewColumn("remark", "Remarks", "備註", multiLanguageTextParser),
	NewColumn("listedDate", "Date Listed", "上市日", dateParser),
	NewColumn("market", "Market", "市場別", multiLanguageTextParser),
	NewColumn("issuedDate", "Date Issued", "發行日", dateParser),
	NewColumn("maturityDate", "Muturity", "到期日", dateParser),
	NewColumn("interestRate", "Interest", "利率值", numberParser),
	NewColumn("registrationDate", "Date of Registering", "登錄日", dateParser),
	NewColumn("announcementDate", "Date Issued", "掛牌日", dateParser),
	NewColumn("securityName", "Security Name", "有價證券名稱", multiLanguageTextParser),
	NewColumn("indexCode,indexName", "Index Code & Index Name", "指數代號及名稱", stringParser, multiLanguageTextParser),
	NewColumn("announcementDate", "Date Announcement", "發布日", dateParser),
	NewColumn("remark", "Remark", "備註", multiLanguageTextParser),
}

type Row struct {
	SecurityCode     *string            `json:"securityCode,omitempty"`
	SecurityName     *MultiLanguageText `json:"securityName,omitempty"`
	IndexCode        *string            `json:"indexCode,omitempty"`
	IndexName        *MultiLanguageText `json:"indexName,omitempty"`
	ISINCode         *string            `json:"isinCode,omitempty"`
	IssuedDate       *civil.Date        `json:"issuedDate,omitempty"`
	ListedDate       *civil.Date        `json:"listedDate,omitempty"`
	MaturityDate     *civil.Date        `json:"maturityDate,omitempty"`
	RegistrationDate *civil.Date        `json:"registrationDate,omitempty"`
	AnnouncementDate *civil.Date        `json:"announcementDate,omitempty"`
	InterestRate     *float64           `json:"interestRate,omitempty"`
	Market           *MultiLanguageText `json:"market,omitempty"`
	IndustrialGroup  *MultiLanguageText `json:"industrialGroup,omitempty"`
	CFICode          *string            `json:"cfiCode,omitempty"`
	Remark           *MultiLanguageText `json:"remark,omitempty"`
}

type Table struct {
	Index       int                `json:"index"`
	Title       MultiLanguageText  `json:"title"`
	Subtitle    *MultiLanguageText `json:"subtitle,omitempty"`
	UpdatedDate civil.Date         `json:"updatedDate"`
	Columns     []Column           `json:"columns"`
	Rows        []Row              `json:"rows"`
}
