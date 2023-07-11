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
		userAdmin := v1.Group("/users").Use(middlewares.IsAdmin())
		{
			userAdmin.GET("/", routes.GetUsers)
			userAdmin.DELETE("/:id", routes.DeleteUser)
		}
		user := v1.Group("/users").Use(middlewares.Auth())
		{
			user.GET("/:id", routes.GetUserByID)
			user.PUT("/:id", routes.PutUser)
		}

		orderAdmin := v1.Group("/orders").Use(middlewares.IsAdmin())
		{
			orderAdmin.GET("/", routes.GetOrders)
		}
		order := v1.Group("/orders").Use(middlewares.Auth())
		{
			order.GET("/:id", routes.GetOrderByID)
			order.POST("/", routes.PostOrder)
			order.PUT("/:id", routes.PutOrder)
			order.DELETE("/:id", routes.DeleteOrder)
		}

		bookAdmin := v1.Group("/books").Use(middlewares.IsAdmin())
		{
			bookAdmin.POST("/", routes.PostBook)
			bookAdmin.PUT("/:id", routes.PutBook)
			bookAdmin.DELETE("/:id", routes.DeleteBook)
		}
		book := v1.Group("/books").Use(middlewares.Auth())
		{
			book.GET("/", routes.GetBooks)
			book.GET("/:id", routes.GetBookByID)
		}

		transactionAdmin := v1.Group("/transactions").Use(middlewares.IsAdmin())
		{
			transactionAdmin.GET("/", routes.GetTransaction)
			transactionAdmin.GET("/book/:id", routes.GetTransactionByBookID)
		}
		transaction := v1.Group("/transactions").Use(middlewares.Auth())
		{
			transaction.GET("/order/:id", routes.GetTransactionByOrderID)
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
