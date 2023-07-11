package routes

import (
	"golang_basic_project/config"
	"golang_basic_project/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type RequestTransaction struct {
	OrderID uint `json:"order_id" binding:"required"`
	BookID  uint `json:"book_id" binding:"required"`
	Price   uint `json:"price" binding:"required"`
}

type ResponseTransaction struct {
	OrderID   uint   `json:"order_id"`
	BookID    uint   `json:"book_id"`
	OrderCode string `json:"order_code"`
	BookName  string `json:"book_name"`
	Price     uint   `json:"price"`
}

func GetTransaction(c *gin.Context) {
	orderBooks := []models.OrderBook{}
	config.DB.Preload(clause.Associations).Find(&orderBooks)
	resTransactions := []ResponseTransaction{}

	for _, ordBook := range orderBooks {
		resTrans := ResponseTransaction{
			OrderID:   ordBook.OrderID,
			BookID:    ordBook.BookID,
			OrderCode: ordBook.Order.OrderCode,
			BookName:  ordBook.Book.Name,
			Price:     ordBook.Price,
		}
		resTransactions = append(resTransactions, resTrans)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success get all transactions",
		"data":    resTransactions,
	})
}

func GetTransactionByOrderID(c *gin.Context) {
	id := c.Param("id")
	var orderBooks []models.OrderBook
	data := config.DB.Preload(clause.Associations).Find(&orderBooks, "order_id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}
	resTransactions := []ResponseTransaction{}
	for _, ordBook := range orderBooks {
		resTrans := ResponseTransaction{
			OrderID:   ordBook.OrderID,
			BookID:    ordBook.BookID,
			OrderCode: ordBook.Order.OrderCode,
			BookName:  ordBook.Book.Name,
			Price:     ordBook.Price,
		}
		resTransactions = append(resTransactions, resTrans)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success get transaction by order id",
		"data":    resTransactions,
	})
}

func GetTransactionByBookID(c *gin.Context) {
	id := c.Param("id")
	var orderBooks []models.OrderBook
	data := config.DB.Preload(clause.Associations).Find(&orderBooks, "book_id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}
	resTransactions := []ResponseTransaction{}
	for _, ordBook := range orderBooks {
		resTrans := ResponseTransaction{
			OrderID:   ordBook.OrderID,
			BookID:    ordBook.BookID,
			OrderCode: ordBook.Order.OrderCode,
			BookName:  ordBook.Book.Name,
			Price:     ordBook.Price,
		}
		resTransactions = append(resTransactions, resTrans)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success get transaction by book id",
		"data":    resTransactions,
	})
}

func PostTransactionOrder(c *gin.Context) {
	var orderBook models.OrderBook
	reqTrans := RequestTransaction{}
	if err := c.ShouldBindJSON(&reqTrans); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"data":    err.Error(),
		})
		c.Abort()
		return
	}
	// check order id exist
	var order models.Order
	checkOrderID := config.DB.First(&order, "id = ?", reqTrans.OrderID)
	if checkOrderID.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Order ID not found",
			"error":   checkOrderID.Error.Error(),
		})
		return
	}
	// check book id exist
	var book models.Book
	checkBookID := config.DB.First(&book, "id = ?", reqTrans.BookID)
	if checkBookID.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Book ID not found",
			"error":   checkBookID.Error.Error(),
		})
		return
	}
	// insert data
	orderBook.OrderID = reqTrans.OrderID
	orderBook.BookID = reqTrans.BookID
	orderBook.Price = reqTrans.Price
	insert := config.DB.Create(&orderBook)
	if insert.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
			"data":    insert.Error.Error(),
		})
		c.Abort()
		return
	}
	// response
	c.JSON(http.StatusOK, gin.H{
		"message": "Success post data transaction",
		"data":    reqTrans,
	})
}
