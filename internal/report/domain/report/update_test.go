package report

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestReportUpdate_LockedRowAllowsPassiveResequence(t *testing.T) {
	expr, _ := NewExpression(uuid.New(), ExpressionNone, nil, nil, nil)
	lockedRow, _ := NewRow(uuid.New(), "LOCKED", "Locked", 1, intPtr(1), true, 1, false, false, false, expr, nil, nil)
	column, _ := NewColumn(uuid.New(), "本月金额", ColumnPeriodAmount, 1)
	r, _ := New(uuid.New(), uuid.New(), uuid.Nil, "Template", true, ClassIncomeStatement, []*Column{column}, []*Row{lockedRow})

	newExpr, _ := NewExpression(uuid.New(), ExpressionNone, nil, nil, nil)
	err := r.UpdateStructure(UpdateParams{
		Rows: []UpdateRowParams{
			{
				RowCode:    "CUSTOM",
				Text:       "Custom",
				LineNo:     intPtr(1),
				ShowLineNo: true,
				SumFactor:  1,
				Expression: newExpr,
			},
			{
				RowId:      lockedRow.Id(),
				RowCode:    "LOCKED",
				Text:       "Locked",
				LineNo:     intPtr(2),
				ShowLineNo: true,
				SumFactor:  1,
				Expression: expr,
			},
		},
	})

	require.NoError(t, err)
	require.Equal(t, 2, r.Rows()[1].Sequence())
}

func TestReportUpdate_LockedRowCannotMoveAcrossExistingRows(t *testing.T) {
	expr, _ := NewExpression(uuid.New(), ExpressionNone, nil, nil, nil)
	beforeRow, _ := NewRow(uuid.New(), "BEFORE", "Before", 1, intPtr(1), true, 1, true, true, false, expr, nil, nil)
	lockedRow, _ := NewRow(uuid.New(), "LOCKED", "Locked", 2, intPtr(2), true, 1, false, false, false, expr, nil, nil)
	afterRow, _ := NewRow(uuid.New(), "AFTER", "After", 3, intPtr(3), true, 1, true, true, false, expr, nil, nil)
	column, _ := NewColumn(uuid.New(), "本月金额", ColumnPeriodAmount, 1)
	r, _ := New(uuid.New(), uuid.New(), uuid.Nil, "Template", true, ClassIncomeStatement, []*Column{column}, []*Row{beforeRow, lockedRow, afterRow})

	err := r.UpdateStructure(UpdateParams{
		Rows: []UpdateRowParams{
			updateParamsFromRow(beforeRow),
			updateParamsFromRow(afterRow),
			updateParamsFromRow(lockedRow),
		},
	})

	require.Error(t, err)
}

func TestReportUpdate_LockedRowCannotEditOrDelete(t *testing.T) {
	expr, _ := NewExpression(uuid.New(), ExpressionNone, nil, nil, nil)
	lockedRow, _ := NewRow(uuid.New(), "LOCKED", "Locked", 1, intPtr(1), true, 1, false, false, false, expr, nil, nil)
	column, _ := NewColumn(uuid.New(), "本月金额", ColumnPeriodAmount, 1)
	r, _ := New(uuid.New(), uuid.New(), uuid.Nil, "Template", true, ClassIncomeStatement, []*Column{column}, []*Row{lockedRow})

	err := r.UpdateStructure(UpdateParams{
		Rows: []UpdateRowParams{
			{
				RowId:      lockedRow.Id(),
				RowCode:    "LOCKED",
				Text:       "Changed",
				LineNo:     intPtr(1),
				ShowLineNo: true,
				SumFactor:  1,
				Expression: expr,
			},
		},
	})
	require.Error(t, err)

	err = r.UpdateStructure(UpdateParams{Rows: []UpdateRowParams{}})
	require.Error(t, err)
}

func intPtr(v int) *int {
	return &v
}

func updateParamsFromRow(row *Row) UpdateRowParams {
	return UpdateRowParams{
		RowId:      row.Id(),
		RowCode:    row.RowCode(),
		Text:       row.Text(),
		LineNo:     row.LineNo(),
		ShowLineNo: row.ShowLineNo(),
		SumFactor:  row.SumFactor(),
		Expression: row.Expression(),
		Rows:       rowsToUpdateParams(row.Rows()),
	}
}

func rowsToUpdateParams(rows []*Row) []UpdateRowParams {
	params := make([]UpdateRowParams, 0, len(rows))
	for _, row := range rows {
		params = append(params, updateParamsFromRow(row))
	}
	return params
}
