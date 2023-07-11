package main

import (
	"golang_basic_project/config"
	"golang_basic_project/middlewares"
	"golang_basic_project/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.Default()
	r.GET("/", getHome)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", routes.RegisterUser)
			auth.POST("/login", routes.GenerateToken)
		}
		user := v1.Group("/users").Use(middlewares.Auth())
		{
			user.GET("/", routes.GetUsers)
			user.GET("/:id", routes.GetUserByID)
			user.PUT("/:id", routes.PutUser)
			user.DELETE("/:id", routes.DeleteUser)
		}
		order := v1.Group("/orders").Use(middlewares.Auth())
		{
			order.GET("/", routes.GetOrders)
			order.GET("/:id", routes.GetOrderByID)
			order.POST("/", routes.PostOrder)
			order.PUT("/:id", routes.PutOrder)
			order.DELETE("/:id", routes.DeleteOrder)
		}
		book := v1.Group("/books").Use(middlewares.Auth())
		{
			book.GET("/", routes.GetBooks)
			book.GET("/:id", routes.GetBookByID)
			book.POST("/", routes.PostBook)
			book.PUT("/:id", routes.PutBook)
			book.DELETE("/:id", routes.DeleteBook)
		}

		transaction := v1.Group("/transactions")
		{
			transaction.GET("/", routes.GetTransaction)
			transaction.GET("/order/:id", routes.GetTransactionByOrderID)
			transaction.GET("/book/:id", routes.GetTransactionByBookID)
			transaction.POST("/order", routes.PostTransactionOrder)
		}
	}

	r.Run()
}

func getHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome",
	})
}
