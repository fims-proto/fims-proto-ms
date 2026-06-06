package app

import (
	"github/fims-proto/fims-proto-ms/internal/dimension/app/command"
	"github/fims-proto/fims-proto-ms/internal/dimension/app/query"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain"
)

type Queries struct {
	PagingCategories query.PagingCategoriesHandler
	CategoryById     query.CategoryByIdHandler
	CategoriesByIds  query.CategoriesByIdsHandler
	PagingOptions    query.PagingOptionsHandler
	OptionsByIds     query.OptionsByIdsHandler
	ValidateOptions  query.ValidateOptionsHandler
}

type Commands struct {
	CreateCategory command.CreateCategoryHandler
	UpdateCategory command.UpdateCategoryHandler
	DeleteCategory command.DeleteCategoryHandler

	CreateOption command.CreateOptionHandler
	UpdateOption command.UpdateOptionHandler
	DeleteOption command.DeleteOptionHandler

	Migrate command.MigrationHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	repo domain.Repository,
	readModel query.DimensionReadModel,
) {
	a.Queries = Queries{
		PagingCategories: query.NewPagingCategoriesHandler(readModel),
		CategoryById:     query.NewCategoryByIdHandler(readModel),
		CategoriesByIds:  query.NewCategoriesByIdsHandler(readModel),
		PagingOptions:    query.NewPagingOptionsHandler(readModel),
		OptionsByIds:     query.NewOptionsByIdsHandler(readModel),
		ValidateOptions:  query.NewValidateOptionsHandler(readModel),
	}
	a.Commands = Commands{
		CreateCategory: command.NewCreateCategoryHandler(repo),
		UpdateCategory: command.NewUpdateCategoryHandler(repo),
		DeleteCategory: command.NewDeleteCategoryHandler(repo),

		CreateOption: command.NewCreateOptionHandler(repo),
		UpdateOption: command.NewUpdateOptionHandler(repo),
		DeleteOption: command.NewDeleteOptionHandler(repo),

		Migrate: command.NewMigrationHandler(repo),
	}
}
