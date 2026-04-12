package report

import (
	commonerrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
)

func (i *Item) UpdateText(text string) error {
	if text == "" {
		return commonerrors.NewInvalidInputError(commonerrors.SlugReportItemEmptyText)
	}

	if text != i.text && !i.isEditable {
		return commonerrors.NewInvalidInputError(commonerrors.SlugReportItemNotEditable)
	}

	i.text = text
	return nil
}

func (i *Item) UpdateSumFactor(sumFactor int) error {
	if sumFactor != -1 && sumFactor != 0 && sumFactor != 1 {
		return commonerrors.NewInvalidInputError(commonerrors.SlugReportItemInvalidSumFactor)
	}

	if sumFactor != i.sumFactor && !i.isEditable {
		return commonerrors.NewInvalidInputError(commonerrors.SlugReportItemNotEditable)
	}

	i.sumFactor = sumFactor
	return nil
}

func (i *Item) UpdateDataSource(dataSource data_source.DataSource, formulas []*Formula) error {
	if dataSource != data_source.Formulas && len(formulas) > 0 {
		return commonerrors.NewInvalidInputError(commonerrors.SlugReportItemInvalidDataSourceWithForms)
	}

	if dataSource != i.dataSource && !i.isEditable {
		return commonerrors.NewInvalidInputError(commonerrors.SlugReportItemNotEditable)
	}

	if !formulasEqual(formulas, i.formulas) && !i.isEditable {
		return commonerrors.NewInvalidInputError(commonerrors.SlugReportItemNotEditable)
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
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}
