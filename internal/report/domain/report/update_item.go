package report

import (
	commonerrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
)

func (i *Item) UpdateText(text string) error {
	if text == "" {
		return commonerrors.NewSlugError("report-item-emptyText")
	}

	i.text = text
	return nil
}

func (i *Item) UpdateSumFactor(sumFactor int) error {
	if sumFactor != -1 && sumFactor != 0 && sumFactor != 1 {
		return commonerrors.NewSlugError("report-item-invalidSumFactor")
	}

	i.sumFactor = sumFactor
	return nil
}

func (i *Item) UpdateDataSource(dataSource data_source.DataSource, formulas []*Formula) error {
	if dataSource != data_source.Formulas && len(formulas) > 0 {
		return commonerrors.NewSlugError("report-item-invalidDataSourceWithFormulas")
	}

	i.dataSource = dataSource
	i.formulas = formulas
	return nil
}
