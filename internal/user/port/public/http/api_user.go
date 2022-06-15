package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/user/app"
	"github/fims-proto/fims-proto-ms/internal/user/app/command"
)

type Handler struct {
	app *app.Application
}

func NewHandler(app *app.Application) Handler {
	if app == nil {
		panic("nil application")
	}
	return Handler{app: app}
}

// ReadUserById godoc
// @Summary Show user by id
// @Description Show user by id
// @Tags users
// @Accept application/json
// @Produce application/json
// @Param userId path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /user/{userId} [get]
func (h Handler) ReadUserById(c *gin.Context) {
	user, err := h.app.Queries.ReadUsers.HandleReadById(c, uuid.MustParse(c.Param("userId")))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if user.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, mapFromUserQuery(user))
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user
// @Tags users
// @Accept application/json
// @Produce application/json
// @Param userId path string true "User ID"
// @Param UpdateUserRequest body UpdateUserRequest true "Update user request"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /user/{userId} [patch]
func (h Handler) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var traits json.RawMessage
	if err := traits.UnmarshalJSON([]byte(req.Traits)); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.UpdateUserCmd{
		Id:     uuid.MustParse(c.Param("userId")),
		Traits: traits,
	}
	if err := h.app.Commands.UpdateUser.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/user/:userId", h.ReadUserById)
	r.PATCH("/user/:userId", h.UpdateUser)
}
