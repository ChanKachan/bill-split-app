package http

import (
	"bill-split/internal/repository"

	"github.com/gin-gonic/gin"
)

type GroupService interface {
}
type group struct {
	groupRepo repository.GroupRepository
}

func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &group{
		groupRepo: groupRepo,
	}
}

// todo: Нужно брать данные АВТОРИЗИРУЕМОГО пользователя (user_id) и записать его как создателя
func (gr *group) CreateGroup(c *gin.Context) {
	//groupData := groupMembers.GroupMembers{}
	//
	//if err := c.ShouldBind(&g); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}

}
