package filterable

import (
	"testing"
)

func TestFilter(t *testing.T) {
	expression := `and(ctn(name, "a"),or(eq(title,"哈哈"),lt(title, 1.3)))`
	filter := &FilterExpr{Buffer: expression}
	filter.Init()
	filter.Print()
	if err := filter.Parse(); err != nil {
		t.Fatal(err)
	}
	// filter.PrintSyntaxTree()
	_, err := filter.ParseAsFilterNode()
	// strSQL, err := assembleSQL(node)
	// println(strSQL)
	if err != nil {
		t.Fatal(err)
	}
}
