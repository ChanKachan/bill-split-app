package handler

import (
	"context"
	"errors"
	"github.com/ChanKachan/bill-split-app/internal/domain/entity/cost"
	"github.com/ChanKachan/bill-split-app/internal/domain/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CostHandler interface {
	CreateCost(c *gin.Context)
	GetCost(c *gin.Context)
	GetGroupCosts(c *gin.Context)
	GetMyCosts(c *gin.Context)
	UpdateCost(c *gin.Context)
	DeleteCost(c *gin.Context)
}

type costHandler struct {
	costService service.CostService
}

func NewCostHandler(costService service.CostService) CostHandler {
	return &costHandler{
		costService: costService,
	}
}

// CreateCost создаёт новый расход
// POST /cost/create
func (h *costHandler) CreateCost(c *gin.Context) {
	var costData cost.Cost
	if err := c.ShouldBindJSON(&costData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем ID текущего пользователя из контекста
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	costData.UserId = userID

	// Валидация
	if costData.GroupId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id is required"})
		return
	}
	if costData.Sum <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sum must be greater than 0"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := h.costService.CreateCost(ctx, costData)
	if err != nil {
		log.Printf("CreateCost | Failed to create cost: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Cost created successfully",
	})
}

// GetCost получает расход по ID
// GET /cost/:id
func (h *costHandler) GetCost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cost id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	costData, err := h.costService.GetCost(ctx, id)
	if err != nil {
		log.Printf("GetCost | Failed to get cost: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Cost not found"})
		return
	}

	// Проверяем доступ пользователя к расходу
	userID, _ := getUserIDFromContext(c)
	if costData.UserId != userID {
		// Можно добавить проверку на принадлежность к группе
	}

	c.JSON(http.StatusOK, costData)
}

// GetGroupCosts получает все расходы группы
// GET /cost/group/:groupId
func (h *costHandler) GetGroupCosts(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	costs, err := h.costService.GetGroupCosts(ctx, groupID)
	if err != nil {
		log.Printf("GetGroupCosts | Failed to get costs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group_id": groupID,
		"costs":    costs,
	})
}

// GetMyCosts получает все расходы текущего пользователя
// GET /cost/my
func (h *costHandler) GetMyCosts(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	costs, err := h.costService.GetUserCosts(ctx, userID)
	if err != nil {
		log.Printf("GetMyCosts | Failed to get costs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"costs":   costs,
	})
}

// UpdateCost обновляет расход
// PUT /cost/:id
func (h *costHandler) UpdateCost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cost id"})
		return
	}

	var costData cost.Cost
	if err := c.ShouldBindJSON(&costData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	costData.Id = id
	costData.UserId = userID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.costService.UpdateCost(ctx, costData)
	if err != nil {
		log.Printf("UpdateCost | Failed to update cost: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cost updated successfully",
	})
}

// DeleteCost удаляет расход
// DELETE /cost/:id
func (h *costHandler) DeleteCost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cost id"})
		return
	}

	// Здесь можно добавить проверку прав на удаление
	// (например, только создатель расхода или админ группы)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.costService.DeleteCost(ctx, id)
	if err != nil {
		log.Printf("DeleteCost | Failed to delete cost: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cost deleted successfully",
	})
}

// Вспомогательная функция для получения userID из контекста
func getUserIDFromContext(c *gin.Context) (int, error) {
	userInfoAny, ok := c.Get("userID")
	if !ok {
		return 0, errors.New("user not authenticated")
	}

	userID, ok := userInfoAny.(int)
	if !ok {
		return 0, errors.New("invalid user ID in context")
	}

	return userID, nil
}
