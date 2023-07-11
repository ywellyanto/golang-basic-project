package routes

import (
	"golang_basic_project/auth"
	"golang_basic_project/config"
	"golang_basic_project/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	var requestRegister models.RequestRegister
	if err := c.ShouldBindJSON(&requestRegister); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}
	user.Username = requestRegister.Username
	user.Email = requestRegister.Email
	user.Password = requestRegister.Password
	user.Role = requestRegister.Role
	// hash password
	err := user.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed hash password",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}
	// insert user to DB
	insert := config.DB.Create(&user)
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
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
	return
}

func GenerateToken(c *gin.Context) {
	request := models.TokenRequest{}
	user := models.User{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}
	// check email
	checkEmail := config.DB.Where("email = ?", request.Email).First(&user)
	if checkEmail.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Email not found",
			"error":   checkEmail.Error.Error(),
		})
		c.Abort()
		return
	}
	// check password
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Password not match",
			"error":   credentialError.Error(),
		})
		c.Abort()
		return
	}
	// generate token
	tokenString, err := auth.GenerateJWT(user.Email, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}
	// response token
	c.JSON(http.StatusCreated, gin.H{
		"token": tokenString,
	})
	return
}

func GetUsers(c *gin.Context) {
	users := []models.User{}
	config.DB.Preload(clause.Associations).Find(&users)

	responseUsers := []models.ResponseUser{}

	for _, u := range users {
		orders := []models.ResponseOrderSingle{}
		for _, o := range u.Orders {
			ord := models.ResponseOrderSingle{
				ID:          o.ID,
				OrderCode:   o.OrderCode,
				Description: o.Description,
				CreatedAt:   o.CreatedAt,
			}

			orders = append(orders, ord)
		}

		user := models.ResponseUser{
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
			UserDetail: models.ResponseUserDetail{
				Name:    u.UserDetail.Name,
				Address: u.UserDetail.Address,
			},
			Orders: orders,
		}

		responseUsers = append(responseUsers, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success get users",
		"data":    responseUsers,
	})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	data := config.DB.Preload(clause.Associations).First(&user, "id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}

	orders := []models.ResponseOrderSingle{}
	for _, o := range user.Orders {
		ord := models.ResponseOrderSingle{
			ID:          o.ID,
			OrderCode:   o.OrderCode,
			Description: o.Description,
			CreatedAt:   o.CreatedAt,
		}

		orders = append(orders, ord)
	}

	responseUser := models.ResponseUser{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		UserDetail: models.ResponseUserDetail{
			Name:    user.UserDetail.Name,
			Address: user.UserDetail.Address,
		},
		Orders: orders,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success get user",
		"data":    responseUser,
	})
}

func PutUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	data := config.DB.Preload(clause.Associations).First(&user, "id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}

	var reqUser models.RequestUser
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	user.Username = reqUser.Username
	user.Role = reqUser.Role
	user.UserDetail.Name = reqUser.Name
	user.UserDetail.Address = reqUser.Address

	update := config.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user)
	if update.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
			"error":   update.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success update user",
		"data":    reqUser,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	data := config.DB.First(&user, "id = ?", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
			"error":   data.Error.Error(),
		})
		return
	}

	config.DB.Delete(&user, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete user success",
	})
}
