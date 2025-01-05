package command

type InitializeCmdReport struct {
	Title       string                 `json:"title"`
	Class       string                 `json:"class"`
	AmountTypes []string               `json:"amountTypes"`
	Sections    []InitializeCmdSection `json:"sections"`
}

type InitializeCmdSection struct {
	Title    string                 `json:"title"`
	Sections []InitializeCmdSection `json:"sections"`
	Items    []InitializeCmdItem    `json:"items"`
}

type InitializeCmdItem struct {
	Text             string                 `json:"text"`
	Level            int                    `json:"level"`
	SumFactor        int                    `json:"sumFactor"`
	DisplaySumFactor bool                   `json:"displaySumFactor"`
	DataSource       string                 `json:"dataSource"`
	Formulas         []InitializeCmdFormula `json:"formulas"`
	IsEditable       bool                   `json:"isEditable"`
	IsBreakdownItem  bool                   `json:"isBreakdownItem"`
	IsAbleToAddChild bool                   `json:"isAbleToAddChild"`
	IsAbleToAddLeaf  bool                   `json:"isAbleToAddLeaf"`
}

type InitializeCmdFormula struct {
	AccountNumber string `json:"accountNumber"`
	SumFactor     int    `json:"sumFactor"`
	Rule          string `json:"rule"`
}
