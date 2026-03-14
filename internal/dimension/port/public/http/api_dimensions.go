package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/dimension/app/command"
	"github/fims-proto/fims-proto-ms/internal/dimension/app/query"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchCategories godoc
//
//	@Summary	Search dimension categories
//	@Tags		dimension
//	@Param		sobId	path		string	true	"Sob ID"
//	@Success	200		{object}	data.PageResponse[CategoryResponse]
//	@Router		/sob/{sobId}/dimension/categories [get]
func (h Handler) SearchCategories(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.DimensionCategory], error) {
			return h.app.Queries.PagingCategories.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		categoryDTOToVO,
	)
}

// ReadCategoryById godoc
//
//	@Summary	Get a dimension category by ID
//	@Tags		dimension
//	@Param		sobId		path		string	true	"Sob ID"
//	@Param		categoryId	path		string	true	"Category ID"
//	@Success	200			{object}	CategoryResponse
//	@Failure	404
//	@Router		/sob/{sobId}/dimension/category/{categoryId} [get]
func (h Handler) ReadCategoryById(c *gin.Context) {
	v, err := h.app.Queries.CategoryById.Handle(c, uuid.MustParse(c.Param("categoryId")))
	if err != nil {
		_ = c.Error(err)
		return
	}

	if v.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, categoryDTOToVO(v))
}

// CreateCategory godoc
//
//	@Summary	Create a dimension category
//	@Tags		dimension
//	@Param		sobId					path	string					true	"Sob ID"
//	@Param		CreateCategoryRequest	body	CreateCategoryRequest	true	"Create category request"
//	@Success	201
//	@Router		/sob/{sobId}/dimension/categories [post]
func (h Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.app.Commands.CreateCategory.Handle(c, command.CreateCategoryCmd{
		CategoryId: uuid.New(),
		SobId:      uuid.MustParse(c.Param("sobId")),
		Name:       req.Name,
	}); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateCategory godoc
//
//	@Summary	Update a dimension category
//	@Tags		dimension
//	@Param		sobId					path	string					true	"Sob ID"
//	@Param		categoryId				path	string					true	"Category ID"
//	@Param		UpdateCategoryRequest	body	UpdateCategoryRequest	true	"Update category request"
//	@Success	200
//	@Router		/sob/{sobId}/dimension/category/{categoryId} [patch]
func (h Handler) UpdateCategory(c *gin.Context) {
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.app.Commands.UpdateCategory.Handle(c, command.UpdateCategoryCmd{
		CategoryId: uuid.MustParse(c.Param("categoryId")),
		NewName:    req.Name,
	}); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

// DeleteCategory godoc
//
//	@Summary	Delete a dimension category
//	@Tags		dimension
//	@Param		sobId		path	string	true	"Sob ID"
//	@Param		categoryId	path	string	true	"Category ID"
//	@Success	204
//	@Router		/sob/{sobId}/dimension/category/{categoryId} [delete]
func (h Handler) DeleteCategory(c *gin.Context) {
	if err := h.app.Commands.DeleteCategory.Handle(c, command.DeleteCategoryCmd{
		CategoryId: uuid.MustParse(c.Param("categoryId")),
	}); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchOptions godoc
//
//	@Summary	Search dimension options within a category
//	@Tags		dimension
//	@Param		sobId		path		string	true	"Sob ID"
//	@Param		categoryId	path		string	true	"Category ID"
//	@Success	200			{object}	data.PageResponse[OptionResponse]
//	@Router		/sob/{sobId}/dimension/category/{categoryId}/options [get]
func (h Handler) SearchOptions(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.DimensionOption], error) {
			return h.app.Queries.PagingOptions.Handle(c, uuid.MustParse(c.Param("categoryId")), pageRequest)
		},
		optionDTOToVO,
	)
}

// CreateOption godoc
//
//	@Summary	Create a dimension option within a category
//	@Tags		dimension
//	@Param		sobId				path	string				true	"Sob ID"
//	@Param		categoryId			path	string				true	"Category ID"
//	@Param		CreateOptionRequest	body	CreateOptionRequest	true	"Create option request"
//	@Success	201
//	@Router		/sob/{sobId}/dimension/category/{categoryId}/options [post]
func (h Handler) CreateOption(c *gin.Context) {
	var req CreateOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.app.Commands.CreateOption.Handle(c, command.CreateOptionCmd{
		OptionId:   uuid.New(),
		CategoryId: uuid.MustParse(c.Param("categoryId")),
		Name:       req.Name,
	}); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateOption godoc
//
//	@Summary	Update a dimension option
//	@Tags		dimension
//	@Param		sobId				path	string				true	"Sob ID"
//	@Param		categoryId			path	string				true	"Category ID"
//	@Param		optionId			path	string				true	"Option ID"
//	@Param		UpdateOptionRequest	body	UpdateOptionRequest	true	"Update option request"
//	@Success	200
//	@Router		/sob/{sobId}/dimension/category/{categoryId}/option/{optionId} [patch]
func (h Handler) UpdateOption(c *gin.Context) {
	var req UpdateOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.app.Commands.UpdateOption.Handle(c, command.UpdateOptionCmd{
		OptionId: uuid.MustParse(c.Param("optionId")),
		NewName:  req.Name,
	}); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

// DeleteOption godoc
//
//	@Summary	Delete a dimension option
//	@Tags		dimension
//	@Param		sobId		path	string	true	"Sob ID"
//	@Param		categoryId	path	string	true	"Category ID"
//	@Param		optionId	path	string	true	"Option ID"
//	@Success	204
//	@Router		/sob/{sobId}/dimension/category/{categoryId}/option/{optionId} [delete]
func (h Handler) DeleteOption(c *gin.Context) {
	if err := h.app.Commands.DeleteOption.Handle(c, command.DeleteOptionCmd{
		OptionId: uuid.MustParse(c.Param("optionId")),
	}); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
