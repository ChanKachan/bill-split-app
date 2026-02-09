package handler

import (
	"bill-split/internal/domain/service/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	authorization Authorization
	user          http.UserHttpService
	group         http.GroupService
}

func NewHandlers(
	userHttpService http.UserHttpService,
	groupService http.GroupService,
) *Handlers {
	return &Handlers{
		//authorization:,
		user:  userHttpService,
		group: groupService,
	}
}

func (h *Handlers) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/user/create", h.user.Create)
	r.POST("/group/create", h.group.CreateGroup)
	r.POST("/group/add/user", h.group.AddUserToGroup)
	//r.GET("/base64", http.GetBaseFile)

	r.Run("0.0.0.0:8000")

	return r
}
