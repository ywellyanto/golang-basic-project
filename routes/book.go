package routes

import (
	"golang_basic_project/config"
	"golang_basic_project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	books := []models.Book{}
	config.DB.Find(&books)

	responseBooks := []models.ResponseBook{}

	for _, b := range books {
		book := models.ResponseBook{
			ID:          b.ID,
			Name:        b.Name,
			Code:        b.Code,
			Description: b.Description,
		}

		responseBooks = append(responseBooks, book)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success get books",
		"data":    responseBooks,
	})
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	data := config.DB.First(&book, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}

	responseBook := models.ResponseBook{
		ID:          book.ID,
		Name:        book.Name,
		Code:        book.Code,
		Description: book.Description,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success get book",
		"data":    responseBook,
	})
}

func PostBook(c *gin.Context) {
	var book models.Book
	var requestBook models.RequestBook
	if err := c.ShouldBindJSON(&requestBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}
	book.Name = requestBook.Name
	book.Code = requestBook.Code
	book.Description = requestBook.Description
	// insert book to DB
	insert := config.DB.Create(&book)
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
		"message": "Success post book",
		"data":    requestBook,
	})
}

func PutBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	data := config.DB.First(&book, "id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}

	var reqBook models.RequestBook
	if err := c.ShouldBindJSON(&reqBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	book.Name = reqBook.Name
	book.Code = reqBook.Code
	book.Description = reqBook.Description

	update := config.DB.Model(&book).Where("id = ?", id).Updates(&book)
	if update.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
			"error":   update.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success update book",
		"data":    reqBook,
	})
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	data := config.DB.First(&book, "id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}

	config.DB.Delete(&book, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete book success",
	})
}
