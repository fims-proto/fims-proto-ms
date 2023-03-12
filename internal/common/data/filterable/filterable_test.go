package filterable

import (
	"testing"
)

func TestCalculator(t *testing.T) {
	expression := "(title eq \"哈哈\") AND (accountNum gte 100 OR accountNum lte 50)"
	calc := &FilterExpr{Buffer: expression}
	calc.Init()
	calc.Print()
	if err := calc.Parse(); err != nil {
		t.Fatal(err)
	}
	bff := ""
	// calc.PrettyPrintSyntaxTree(bff)
	calc.PrintSyntaxTree()
	t.Log(bff)
}
