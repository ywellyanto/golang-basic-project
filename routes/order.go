package routes

import (
	"golang_basic_project/config"
	"golang_basic_project/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetOrders(c *gin.Context) {
	orders := []models.Order{}
	config.DB.Preload(clause.Associations).Find(&orders)

	responseOrders := []models.ResponseOrder{}

	for _, o := range orders {
		order := models.ResponseOrder{
			ID:          o.ID,
			OrderCode:   o.OrderCode,
			Description: o.Description,
			CreatedAt:   o.CreatedAt,
			User: models.ResponseUserSingle{
				Username: o.User.Username,
				Email:    o.User.Email,
				Role:     o.User.Role,
			},
		}

		responseOrders = append(responseOrders, order)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success get orders",
		"data":    responseOrders,
	})
}

func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	data := config.DB.Preload(clause.Associations).First(&order, "id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}

	responseOrder := models.ResponseOrder{
		ID:          order.ID,
		OrderCode:   order.OrderCode,
		Description: order.Description,
		CreatedAt:   order.CreatedAt,
		User: models.ResponseUserSingle{
			Username: order.User.Username,
			Email:    order.User.Email,
			Role:     order.User.Role,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success get order",
		"data":    responseOrder,
	})
}

func PostOrder(c *gin.Context) {
	var order models.Order
	var requestOrder models.RequestOrder
	if err := c.ShouldBindJSON(&requestOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}
	order.OrderCode = requestOrder.OrderCode
	order.Description = requestOrder.Description
	order.UserID = requestOrder.UserID
	// insert order to DB
	insert := config.DB.Create(&order)
	if insert.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
			"error":   insert.Error.Error(),
		})
		c.Abort()
		return
	}
	// response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Success post order",
		"data":    requestOrder,
	})
}

func PutOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	data := config.DB.First(&order, "id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}

	var reqOrder models.RequestOrder
	if err := c.ShouldBindJSON(&reqOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	order.OrderCode = reqOrder.OrderCode
	order.Description = reqOrder.Description
	order.UserID = reqOrder.UserID

	update := config.DB.Model(&order).Where("id = ?", id).Updates(&order)
	if update.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
			"error":   update.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success update order",
		"data":    reqOrder,
	})
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	data := config.DB.First(&order, "id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}

	config.DB.Delete(&order, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete order success",
	})
}
