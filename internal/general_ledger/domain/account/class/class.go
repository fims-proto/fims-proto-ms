package class

import "fmt"

type Class int

const (
	Assets           Class = 1 // 资产
	Liabilities      Class = 2 // 负债
	Equities         Class = 3 // 权益
	Costs            Class = 4 // 成本
	ProfitsAndLosses Class = 5 // 损益
	Commons          Class = 7 // 共同
)

var classToString = map[Class]string{
	Assets:           "assets",
	Liabilities:      "liabilities",
	Equities:         "equities",
	Costs:            "costs",
	ProfitsAndLosses: "profits and losses",
	Commons:          "commons",
}

func (c Class) String() string {
	s, ok := classToString[c]
	if ok {
		return s
	}
	return fmt.Sprintf("unknown class [%d]", c)
}
