package item_type

import "fmt"

type ItemType struct {
	slug string
}

func (it ItemType) String() string {
	return it.slug
}

var (
	None            = ItemType{""}
	GrossProfit     = ItemType{"gross_profit"}
	OperatingProfit = ItemType{"operating_profit"}
	TotalProfit     = ItemType{"total_profit"}
	NetProfit       = ItemType{"net_profit"}
)

var stringToItemType = map[string]ItemType{
	"":                 None,
	"gross_profit":     GrossProfit,
	"operating_profit": OperatingProfit,
	"total_profit":     TotalProfit,
	"net_profit":       NetProfit,
}

func FromString(s string) (ItemType, error) {
	itemType, ok := stringToItemType[s]
	if ok {
		return itemType, nil
	}
	return None, fmt.Errorf("unknown item type: %s", s)
}
