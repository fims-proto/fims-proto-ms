package report

import (
	commonerrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
)

func (i *Item) UpdateText(text string) error {
	if text == "" {
		return commonerrors.NewSlugError("report-item-emptyText")
	}

	if text != i.text && !i.isEditable {
		return commonerrors.NewSlugError("report-item-notEditable")
	}

	i.text = text
	return nil
}

func (i *Item) UpdateSumFactor(sumFactor int) error {
	if sumFactor != -1 && sumFactor != 0 && sumFactor != 1 {
		return commonerrors.NewSlugError("report-item-invalidSumFactor")
	}

	if sumFactor != i.sumFactor && !i.isEditable {
		return commonerrors.NewSlugError("report-item-notEditable")
	}

	i.sumFactor = sumFactor
	return nil
}

func (i *Item) UpdateDataSource(dataSource data_source.DataSource, formulas []*Formula) error {
	if dataSource != data_source.Formulas && len(formulas) > 0 {
		return commonerrors.NewSlugError("report-item-invalidDataSourceWithFormulas")
	}

	if dataSource != i.dataSource && !i.isEditable {
		return commonerrors.NewSlugError("report-item-notEditable")
	}

	if !formulasEqual(formulas, i.formulas) && !i.isEditable {
		return commonerrors.NewSlugError("report-item-notEditable")
	}

	i.dataSource = dataSource
	i.formulas = formulas
	return nil
}

func formulasEqual(a, b []*Formula) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
