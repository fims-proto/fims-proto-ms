package command

type InitializeCmdReport struct {
	Title       string                 `json:"title"`
	Class       string                 `json:"class"`
	AmountTypes []string               `json:"amountTypes"`
	Sections    []InitializeCmdSection `json:"sections"`
}

type InitializeCmdSection struct {
	Title       string                 `json:"title"`
	SectionType string                 `json:"sectionType"`
	Sections    []InitializeCmdSection `json:"sections"`
	Items       []InitializeCmdItem    `json:"items"`
}

type InitializeCmdItem struct {
	Text             string                 `json:"text"`
	Level            int                    `json:"level"`
	ItemType         string                 `json:"itemType"`
	SumFactor        int                    `json:"sumFactor"`
	DisplaySumFactor bool                   `json:"displaySumFactor"`
	DataSource       string                 `json:"dataSource"`
	Formulas         []InitializeCmdFormula `json:"formulas"`
	IsEditable       bool                   `json:"isEditable"`
	IsBreakdownItem  bool                   `json:"isBreakdownItem"`
	IsAbleToAddChild bool                   `json:"isAbleToAddChild"`
}

type InitializeCmdFormula struct {
	AccountNumber string `json:"accountNumber"`
	SumFactor     int    `json:"sumFactor"`
	Rule          string `json:"rule"`
}
