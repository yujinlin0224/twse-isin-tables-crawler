package twseisintablescrawler

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

var SupportedColumns = []Column{
	NewColumn("securityCode,securityName", "Security Code & Security Name", "有價證券代號及名稱", StringParser, MultiLanguageTextParser),
	NewColumn("isinCode", "ISIN Code", "國際證券辨識號碼(ISIN Code)", StringParser),
	NewColumn("issuedDate", "Date Issued", "公開發行日", DateParser),
	NewColumn("industrialGroup", "Industrial Group", "產業別", MultiLanguageTextParser),
	NewColumn("cfiCode", "CFICode", "CFICode", StringParser),
	NewColumn("remark", "Remarks", "備註", MultiLanguageTextParser),
	NewColumn("listedDate", "Date Listed", "上市日", DateParser),
	NewColumn("market", "Market", "市場別", MultiLanguageTextParser),
	NewColumn("issuedDate", "Date Issued", "發行日", DateParser),
	NewColumn("inceptionDate", "Date Issued", "基金成立日", DateParser),
	NewColumn("maturityDate", "Muturity", "到期日", DateParser),
	NewColumn("interestRate", "Interest", "利率值", NumberParser),
	NewColumn("registrationDate", "Date of Registering", "登錄日", DateParser),
	NewColumn("announcementDate", "Date Issued", "掛牌日", DateParser),
	NewColumn("securityName", "Security Name", "有價證券名稱", MultiLanguageTextParser),
	NewColumn("indexCode,indexName", "Index Code & Index Name", "指數代號及名稱", StringParser, MultiLanguageTextParser),
	NewColumn("announcementDate", "Date Announcement", "發布日", DateParser),
	NewColumn("remark", "Remark", "備註", MultiLanguageTextParser),
}

type Row struct {
	SecurityCode     *string            `json:"securityCode,omitempty"`
	SecurityName     *MultiLanguageText `json:"securityName,omitempty"`
	IndexCode        *string            `json:"indexCode,omitempty"`
	IndexName        *MultiLanguageText `json:"indexName,omitempty"`
	ISINCode         *string            `json:"isinCode,omitempty"`
	IssuedDate       *civil.Date        `json:"issuedDate,omitempty"`
	InceptionDate    *civil.Date        `json:"inceptionDate,omitempty"`
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
