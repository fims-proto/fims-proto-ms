package filterable

import (
	"testing"
)

func TestFilter(t *testing.T) {
	expression := `and(btw(name, 1, 2),or(eq(title,"哈哈"),lt(title, 1.3)))`
	filter := &FilterExpr{Buffer: expression}
	filter.Init()
	filter.Print()
	if err := filter.Parse(); err != nil {
		t.Fatal(err)
	}
	filter.PrintSyntaxTree()
	_, err := filter.ParseAsFilterNode()
	if err != nil {
		t.Fatal(err)
	}
}
