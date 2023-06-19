package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"server-v2/utils"
	"strconv"
)

func RegisterUser(c *gin.Context) {
	var input models.UserInput
	var detail models.Detail

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Cek keberadaan nama pengguna
	var existingUsername models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUsername).Error; err == nil {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	// Cek keberadaan alamat email
	var existingEmail models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingEmail).Error; err == nil {
		c.JSON(400, gin.H{"error": "Email already exists"})
		return
	}

	detail = models.Detail{
		Balance: 0,
		Credit:  0,
	}
	if err := config.DB.Create(&detail).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	//hash, err := utils.HashPassword(input.Password)
	//if err != nil {
	//	c.JSON(500, gin.H{"error": err.Error()})
	//	return
	//}
	user := models.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  input.Password,
		RoleId:    3,
		DetailId:  detail.Id,
		CompanyID: input.CompanyID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := models.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}

	c.JSON(200, gin.H{"message": "User registered successfully", "user": response})
}

func LoginUser(c *gin.Context) {
	var inputAuth models.Authentication

	if err := c.ShouldBindJSON(&inputAuth); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("Email = ?", inputAuth.Email).Preload("Role").Find(&user).Error; err != nil {
		c.JSON(500, gin.H{"message": "username or password is incorrect"})
		return
	}

	//if !utils.CheckPasswordHash(inputAuth.Password, user.Password) {
	//	c.JSON(401, gin.H{"message": "username or password is incorrect"})
	//	return
	//}

	if inputAuth.Password != user.Password {
		c.JSON(401, gin.H{"message": "username or password is incorrect"})
		return
	}

	//make user.id to string
	id := strconv.Itoa(user.Id)
	token, err := utils.GenerateJWT(user.Email, user.Role.Role, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := models.Token{
		TokenString: token,
		Email:       user.Email,
		Role:        user.Role.Role,
		Name:        user.Username,
	}

	c.JSON(200, gin.H{"message": "Login success", "token": response})

}
