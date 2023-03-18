package filterable

import (
	"testing"
)

func TestFilter(t *testing.T) {
	expression := "in(accountNumber,\"1821\",\"1001\")"
	filter := &FilterExpr{Buffer: expression}
	filter.Init()
	filter.Print()
	if err := filter.Parse(); err != nil {
		t.Fatal(err)
	}
	// filter.PrintSyntaxTree()
	_, err := filter.ParseAsFilterable()
	// _, err = assembleSQL(node)
	// println(strSQL)
	if err != nil {
		t.Fatal(err)
	}
}
