package handler

import (
	"bill-split/internal/domain/entity/transfer"
	"bill-split/internal/domain/service"
	"github.com/gin-gonic/gin"
	"log"
)

type OptimizationHandler interface {
	Optimize(c *gin.Context)
}

type optimizationHandler struct {
	OptimizationService service.OptimizationService
}

func NewOptimizationHandler(OptimizationService service.OptimizationService) OptimizationHandler {
	return &optimizationHandler{
		OptimizationService: OptimizationService,
	}
}

// Метод реализует жадный алгорим.
// Самый большой должник погашает долг самого большого кредитора
func (o *optimizationHandler) Optimize(c *gin.Context) {
	var participants []transfer.Participant

	if err := c.ShouldBind(&participants); err != nil {
		log.Println("OptimizationHandler.Optimize err:", err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	transaction := o.OptimizationService.OptimizeDebts(participants)
	c.JSON(200, gin.H{
		"code": 200,
		"data": transaction,
	})
	return
}
