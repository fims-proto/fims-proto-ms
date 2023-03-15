package filterable

import (
	"strconv"

	"github.com/pkg/errors"
)

type FilterNodeType int

const (
	TypeATOM FilterNodeType = 1 << iota // identity no child filters
	TypeAND                             // and
	TypeOR                              // or
	TypeNOT                             // not, only one child filter
)

type FilterNode interface {
	IsFiltered() bool
	Children() []FilterNode // child filters
	Type() FilterNodeType
}

type filterNodeImpl struct {
	filterableType FilterNodeType
	children       []FilterNode
}

// new

func NewFilterNode(fType FilterNodeType, filters ...FilterNode) FilterNode {
	return &filterNodeImpl{filterableType: fType, children: filters}
}

// impl

func (f *filterNodeImpl) IsFiltered() bool {
	return len(f.children) > 0
}

func (f *filterNodeImpl) Children() []FilterNode {
	return f.children
}

func (f *filterNodeImpl) Type() FilterNodeType {
	return f.filterableType
}

func (fe *FilterExpr) ParseAsFilterNode() (FilterNode, error) {
	node := fe.AST()
	node = node.up
	if node.pegRule != ruleExpr {
		return nil, errors.Errorf("Parse failed at %s", node.String())
	}

	return fe.ParseExpr(node)
}

func (fe *FilterExpr) ParseExpr(node *node32) (FilterNode, error) {
	node = node.up
	if node == nil {
		return nil, errors.Errorf("Pase failed at %s", node.String())
	}
	switch node.pegRule {
	case ruleAndExpr:
		{
			return fe.ParseAndExpr(node)
		}
	case ruleOrExpr:
		{
			return fe.ParseOrExpr(node)
		}
	case ruleNotExpr:
		{
			return fe.ParseNotExpr(node)
		}
	case ruleAtomExpr:
		{
			return fe.ParseAtomExpr(node)
		}
	default:
		{
			return nil, errors.New("unknow rule type")
		}
	}
}

func (fe *FilterExpr) ParseAndExpr(node *node32) (FilterNode, error) {
	var children []FilterNode
	node = node.up.next.next.up
	for node != nil {
		if node.pegRule == ruleExpr {
			f, err := fe.ParseExpr(node)
			if err != nil {
				return nil, err
			}
			children = append(children, f)
		}
		node = node.next
	}
	if children != nil {
		return NewFilterNode(TypeAND, children...), nil
	}
	return nil, errors.Errorf("no child for andExpr")
}

func (fe *FilterExpr) ParseOrExpr(node *node32) (FilterNode, error) {
	var children []FilterNode
	node = node.up.next.next.up
	for node != nil {
		if node.pegRule == ruleExpr {
			f, err := fe.ParseExpr(node)
			if err != nil {
				return nil, err
			}
			children = append(children, f)
		}
		node = node.next
	}
	if children != nil {
		return NewFilterNode(TypeOR, children...), nil
	}
	return nil, errors.Errorf("no child for orExpr")
}

func (fe *FilterExpr) ParseNotExpr(node *node32) (FilterNode, error) {
	var children []FilterNode
	node = node.up.next.next
	f, err := fe.ParseExpr(node)
	if err != nil {
		return nil, err
	}
	children = append(children, f)
	return NewFilterNode(TypeNOT, children...), nil
}

func (fe *FilterExpr) ParseAtomExpr(node *node32) (FilterNode, error) {
	// return types are all filterImpl
	node = node.up
	var op Operator
	switch node.pegRule {
	case ruleEqExpr:
		{
			op = OptEq
		}
	case ruleLtExpr:
		{
			op = OptLt
		}
	case ruleLteExpr:
		{
			op = OptLte
		}
	case ruleGtExpr:
		{
			op = OptGt
		}
	case ruleGteExpr:
		{
			op = OptGte
		}
	case ruleBtwExpr:
		{
			op = OptBtw
		}
	default:
		{
			return nil, errors.New("unknow ruleType")
		}
	}
	field := node.up.next.next
	fieldName := string(fe.buffer[field.begin:field.end])
	filter, err := NewFilter1(fieldName, op, fe.ParseLiterals(field.next.next)...)
	if err != nil {
		return nil, err
	}
	fImpl := filter.(filterImpl)
	if err != nil {
		return fImpl, err
	}
	return nil, err
}

func (fe *FilterExpr) ParseLiterals(node *node32) []any {
	var values []any
	for node != nil {
		if node.pegRule == ruleLiteral {
			literalNode := node.up
			if literalNode.pegRule == ruleStringLiteral {
				values = append(values, string(fe.buffer[literalNode.begin+1:literalNode.end-1]))
			} else if literalNode.pegRule == ruleNumLiteral {
				strVal := string(fe.buffer[literalNode.begin:literalNode.end])
				val, err := strconv.ParseFloat(strVal, 64)
				if err != nil {
					return nil
				}
				values = append(values, val)
			}
		}
		node = node.next
	}
	return values
}
