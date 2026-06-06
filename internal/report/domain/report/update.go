package report

import (
	"fmt"

	"github.com/google/uuid"
)

type UpdateParams struct {
	Title *string
	Rows  []UpdateRowParams
}

type UpdateRowParams struct {
	RowId            uuid.UUID
	RowCode          string
	Text             string
	LineNo           *int
	ShowLineNo       bool
	SumFactor        int
	DisplaySumFactor bool
	Indent           int
	CanEdit          *bool
	CanMove          *bool
	CanAddChild      *bool
	Expression       *Expression
	Rows             []UpdateRowParams
}

func rebuildRows(current []*Row, desired []UpdateRowParams) ([]*Row, error) {
	currentById := make(map[uuid.UUID]*Row)
	for _, row := range current {
		currentById[row.id] = row
	}

	desiredIds := make(map[uuid.UUID]struct{})
	for _, row := range desired {
		if row.RowId != uuid.Nil {
			desiredIds[row.RowId] = struct{}{}
		}
	}
	for _, row := range current {
		if _, ok := desiredIds[row.id]; !ok && !row.canEdit {
			return nil, fmt.Errorf("row %s cannot be deleted", row.rowCode)
		}
	}

	rows := make([]*Row, 0, len(desired))
	for i, desiredRow := range desired {
		rowId := desiredRow.RowId
		if rowId == uuid.Nil {
			rowId = uuid.New()
		}

		canEdit := true
		canMove := true
		canAddChild := false
		if desiredRow.CanEdit != nil {
			canEdit = *desiredRow.CanEdit
		}
		if desiredRow.CanMove != nil {
			canMove = *desiredRow.CanMove
		}
		if desiredRow.CanAddChild != nil {
			canAddChild = *desiredRow.CanAddChild
		}

		if desiredRow.Expression == nil {
			expr, _ := NewExpression(uuid.New(), ExpressionNone, nil, nil, nil)
			desiredRow.Expression = expr
		}

		if desiredRow.RowId != uuid.Nil {
			old := currentById[desiredRow.RowId]
			if old == nil {
				return nil, fmt.Errorf("row %s not found", desiredRow.RowId)
			}
			if !old.canEdit && (old.rowCode != desiredRow.RowCode || old.text != desiredRow.Text || !expressionsEqual(old.expression, desiredRow.Expression)) {
				return nil, fmt.Errorf("row %s cannot be edited", old.rowCode)
			}
			if !old.canMove && lockedRowRelativeOrderChanged(current, desired, old.id) {
				return nil, fmt.Errorf("row %s cannot be moved", old.rowCode)
			}
			if !old.canAddChild && hasAddedChildren(old.rows, desiredRow.Rows) {
				return nil, fmt.Errorf("row %s cannot add child rows", old.rowCode)
			}
			canEdit = old.canEdit
			canMove = old.canMove
			canAddChild = old.canAddChild
		}

		var oldRows []*Row
		if desiredRow.RowId != uuid.Nil {
			oldRows = currentById[desiredRow.RowId].rows
		}
		childRows, err := rebuildRows(oldRows, desiredRow.Rows)
		if err != nil {
			return nil, err
		}

		row, err := NewRow(
			rowId,
			desiredRow.RowCode,
			desiredRow.Text,
			i+1,
			desiredRow.LineNo,
			desiredRow.ShowLineNo,
			desiredRow.SumFactor,
			desiredRow.DisplaySumFactor,
			desiredRow.Indent,
			canEdit,
			canMove,
			canAddChild,
			desiredRow.Expression,
			childRows,
			nil,
		)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func hasAddedChildren(current []*Row, desired []UpdateRowParams) bool {
	currentById := make(map[uuid.UUID]struct{})
	for _, row := range current {
		currentById[row.id] = struct{}{}
	}
	for _, row := range desired {
		if row.RowId == uuid.Nil {
			return true
		}
		if _, ok := currentById[row.RowId]; !ok {
			return true
		}
	}
	return false
}

func lockedRowRelativeOrderChanged(current []*Row, desired []UpdateRowParams, rowId uuid.UUID) bool {
	desiredPosition := make(map[uuid.UUID]int)
	for i, row := range desired {
		if row.RowId != uuid.Nil {
			desiredPosition[row.RowId] = i
		}
	}

	currentPosition := -1
	for i, row := range current {
		if row.id == rowId {
			currentPosition = i
			break
		}
	}
	lockedDesiredPosition, ok := desiredPosition[rowId]
	if currentPosition == -1 || !ok {
		return true
	}

	for i, row := range current {
		if row.id == rowId {
			continue
		}
		otherDesiredPosition, exists := desiredPosition[row.id]
		if !exists {
			continue
		}
		wasBefore := i < currentPosition
		isBefore := otherDesiredPosition < lockedDesiredPosition
		if wasBefore != isBefore {
			return true
		}
	}
	return false
}
