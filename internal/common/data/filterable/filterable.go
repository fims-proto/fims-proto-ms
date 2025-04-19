package filterable

import (
	"fmt"
	"strconv"
)

type Type int

const (
	TypeATOM    Type = 1 << iota // identity one child filters
	TypeAND                      // and
	TypeOR                       // or
	TypeNOT                      // not, only one child filter
	TypeNONE                     // empty not a filter
	TypeRequest                  // current filterable is pageRequestImpl, should not go into assembleSQL
)

type Filterable interface {
	IsFiltered() bool
	Children() []Filterable // child filters
	FilterableType() Type
}

type filterableImpl struct {
	filterableType Type
	children       []Filterable
}

// new

func Unfiltered() Filterable {
	return &filterableImpl{filterableType: TypeNONE, children: nil}
}

func NewFilterable(fType Type, filters ...Filterable) Filterable {
	return &filterableImpl{filterableType: fType, children: filters}
}

func NewFilterableAtom(filter Filter) Filterable {
	return filter.(filterImpl)
}

// impl

func (f *filterableImpl) IsFiltered() bool {
	return len(f.children) > 0
}

func (f *filterableImpl) Children() []Filterable {
	return f.children
}

func (f *filterableImpl) FilterableType() Type {
	return f.filterableType
}

func (fe *FilterAST) ParseAsFilterable() (Filterable, error) {
	node := fe.AST()
	node = node.up
	if node.pegRule != ruleExpr {
		return nil, fmt.Errorf("parse failed at %s", node.String())
	}

	return fe.ParseExpr(node)
}

func (fe *FilterAST) ParseExpr(node *node32) (Filterable, error) {
	if node.up == nil {
		return nil, fmt.Errorf("parse failed at %s", node.String())
	}
	node = node.up
	switch node.pegRule {
	case ruleAndExpr:
		return fe.ParseAndExpr(node)
	case ruleOrExpr:
		return fe.ParseOrExpr(node)
	case ruleNotExpr:
		return fe.ParseNotExpr(node)
	case ruleAtomExpr:
		return fe.ParseAtomExpr(node)
	default:
		return nil, fmt.Errorf("unknow rule type")
	}
}

func (fe *FilterAST) ParseAndExpr(node *node32) (Filterable, error) {
	var children []Filterable
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
		return NewFilterable(TypeAND, children...), nil
	}
	return nil, fmt.Errorf("no child for andExpr")
}

func (fe *FilterAST) ParseOrExpr(node *node32) (Filterable, error) {
	var children []Filterable
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
		return NewFilterable(TypeOR, children...), nil
	}
	return nil, fmt.Errorf("no child for orExpr")
}

func (fe *FilterAST) ParseNotExpr(node *node32) (Filterable, error) {
	var children []Filterable
	node = node.up.next.next
	f, err := fe.ParseExpr(node)
	if err != nil {
		return nil, err
	}
	children = append(children, f)
	return NewFilterable(TypeNOT, children...), nil
}

func (fe *FilterAST) ParseAtomExpr(node *node32) (Filterable, error) {
	// return types are all filterImpl
	node = node.up
	var op Operator
	switch node.pegRule {
	case ruleEqExpr:
		op = OptEq
	case ruleLtExpr:
		op = OptLt
	case ruleLteExpr:
		op = OptLte
	case ruleGtExpr:
		op = OptGt
	case ruleGteExpr:
		op = OptGte
	case ruleBtwExpr:
		op = OptBtw
	case ruleStwExpr:
		op = OptStw
	case ruleCtnExpr:
		op = OptCtn
	case ruleInExpr:
		op = OptIn
	default:
		return nil, fmt.Errorf("unknown ruleType")
	}
	field := node.up.next.next
	fieldName := string(fe.buffer[field.begin:field.end])
	filter, err := NewFilter(fieldName, op, fe.ParseLiterals(field.next.next)...)
	if err != nil {
		return nil, err
	}
	fImpl := filter.(filterImpl)
	return fImpl, err
}

func (fe *FilterAST) ParseLiterals(node *node32) []any {
	var values []any
	if node.pegRule == ruleLiteralList {
		node = node.up
	}
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
