package report

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestReportUpdate_LockedRowAllowsPassiveResequence(t *testing.T) {
	expr, _ := NewExpression(uuid.New(), ExpressionNone, nil, nil, nil)
	lockedRow, _ := NewRow(uuid.New(), "LOCKED", "Locked", 1, intPtr(1), true, 1, false, 0, false, false, false, expr, nil, nil)
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

func TestReportUpdate_PreservesColumnsWhenUpdatingStructure(t *testing.T) {
	expr, _ := NewExpression(uuid.New(), ExpressionNone, nil, nil, nil)
	row, _ := NewRow(uuid.New(), "CUSTOM", "Original", 1, intPtr(1), true, 1, false, 0, true, true, false, expr, nil, nil)
	periodColumn, _ := NewColumn(uuid.New(), "本月金额", ColumnPeriodAmount, 1)
	ytdColumn, _ := NewColumn(uuid.New(), "本年累计金额", ColumnYearToDateAmount, 2)
	r, _ := New(uuid.New(), uuid.New(), uuid.Nil, "Template", true, ClassIncomeStatement, []*Column{periodColumn, ytdColumn}, []*Row{row})

	newTitle := "Updated Template"
	err := r.UpdateStructure(UpdateParams{
		Title: &newTitle,
		Rows: []UpdateRowParams{
			{
				RowId:      row.Id(),
				RowCode:    "CUSTOM",
				Text:       "Updated",
				LineNo:     intPtr(1),
				ShowLineNo: true,
				SumFactor:  1,
				Expression: expr,
			},
		},
	})

	require.NoError(t, err)
	require.Equal(t, newTitle, r.Title())
	require.Equal(t, "Updated", r.Rows()[0].Text())
	require.Len(t, r.Columns(), 2)
	require.Equal(t, periodColumn.Id(), r.Columns()[0].Id())
	require.Equal(t, periodColumn.Label(), r.Columns()[0].Label())
	require.Equal(t, periodColumn.ValueType(), r.Columns()[0].ValueType())
	require.Equal(t, periodColumn.Sequence(), r.Columns()[0].Sequence())
	require.Equal(t, ytdColumn.Id(), r.Columns()[1].Id())
	require.Equal(t, ytdColumn.Label(), r.Columns()[1].Label())
	require.Equal(t, ytdColumn.ValueType(), r.Columns()[1].ValueType())
	require.Equal(t, ytdColumn.Sequence(), r.Columns()[1].Sequence())
}

func TestReportUpdate_DisplaySumFactorSurvivesUpdateAndInstantiate(t *testing.T) {
	expr, _ := NewExpression(uuid.New(), ExpressionNone, nil, nil, nil)
	row, _ := NewRow(uuid.New(), "CUSTOM", "Original", 1, intPtr(1), true, -1, false, 1, true, true, false, expr, nil, nil)
	column, _ := NewColumn(uuid.New(), "本月金额", ColumnPeriodAmount, 1)
	r, _ := New(uuid.New(), uuid.New(), uuid.Nil, "Template", true, ClassIncomeStatement, []*Column{column}, []*Row{row})

	err := r.UpdateStructure(UpdateParams{
		Rows: []UpdateRowParams{
			{
				RowId:            row.Id(),
				RowCode:          "CUSTOM",
				Text:             "Updated",
				LineNo:           intPtr(1),
				ShowLineNo:       true,
				SumFactor:        -1,
				DisplaySumFactor: true,
				Indent:           1,
				Expression:       expr,
			},
		},
	})

	require.NoError(t, err)
	require.True(t, r.Rows()[0].DisplaySumFactor())

	instance, err := r.Instantiate(uuid.New(), uuid.New(), "")

	require.NoError(t, err)
	require.True(t, instance.Rows()[0].DisplaySumFactor())
}

func TestReportUpdate_LockedRowCannotMoveAcrossExistingRows(t *testing.T) {
	expr, _ := NewExpression(uuid.New(), ExpressionNone, nil, nil, nil)
	beforeRow, _ := NewRow(uuid.New(), "BEFORE", "Before", 1, intPtr(1), true, 1, false, 0, true, true, false, expr, nil, nil)
	lockedRow, _ := NewRow(uuid.New(), "LOCKED", "Locked", 2, intPtr(2), true, 1, false, 0, false, false, false, expr, nil, nil)
	afterRow, _ := NewRow(uuid.New(), "AFTER", "After", 3, intPtr(3), true, 1, false, 0, true, true, false, expr, nil, nil)
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
	lockedRow, _ := NewRow(uuid.New(), "LOCKED", "Locked", 1, intPtr(1), true, 1, false, 0, false, false, false, expr, nil, nil)
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

func TestReportUpdate_LockedRowAllowsEquivalentEmptyExpressionReferences(t *testing.T) {
	oldExpr, _ := NewExpression(uuid.New(), ExpressionChildrenSum, nil, nil, nil)
	desiredExpr, _ := NewExpression(
		uuid.New(),
		ExpressionChildrenSum,
		[]LedgerAccountReference{},
		[]CashFlowItemReference{},
		[]RowReference{},
	)
	lockedRow, _ := NewRow(uuid.New(), "LOCKED", "Locked", 1, intPtr(1), true, 1, false, 0, false, false, false, oldExpr, nil, nil)
	column, _ := NewColumn(uuid.New(), "本月金额", ColumnPeriodAmount, 1)
	r, _ := New(uuid.New(), uuid.New(), uuid.Nil, "Template", true, ClassIncomeStatement, []*Column{column}, []*Row{lockedRow})

	err := r.UpdateStructure(UpdateParams{
		Rows: []UpdateRowParams{
			{
				RowId:      lockedRow.Id(),
				RowCode:    "LOCKED",
				Text:       "Locked",
				LineNo:     intPtr(1),
				ShowLineNo: true,
				SumFactor:  1,
				Expression: desiredExpr,
			},
		},
	})

	require.NoError(t, err)
}

func intPtr(v int) *int {
	return &v
}

func updateParamsFromRow(row *Row) UpdateRowParams {
	return UpdateRowParams{
		RowId:            row.Id(),
		RowCode:          row.RowCode(),
		Text:             row.Text(),
		LineNo:           row.LineNo(),
		ShowLineNo:       row.ShowLineNo(),
		SumFactor:        row.SumFactor(),
		DisplaySumFactor: row.DisplaySumFactor(),
		Expression:       row.Expression(),
		Rows:             rowsToUpdateParams(row.Rows()),
	}
}

func rowsToUpdateParams(rows []*Row) []UpdateRowParams {
	params := make([]UpdateRowParams, 0, len(rows))
	for _, row := range rows {
		params = append(params, updateParamsFromRow(row))
	}
	return params
}
