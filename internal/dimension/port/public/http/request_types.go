package http

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateOptionRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateOptionRequest struct {
	Name string `json:"name" binding:"required"`
}
