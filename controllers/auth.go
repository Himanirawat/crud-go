package controllers

import (
	"crud-go/models"
	"crud-go/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// crud operations for mysql
func Register(c *gin.Context) {

	var input utils.RegisterType
	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("%s", err)

		return
	}
	fmt.Print(input)
	u := models.User{}
	u.Username = input.Username
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = string(hashedPassword)
	u.City = input.City
	u.Name = input.Name

	db := models.Connection()
	u.Save(db)
	c.JSON(http.StatusOK, gin.H{
		"message": "registered",
	})
	db.Close()
}

func Login(c *gin.Context) {
	var input utils.RegisterType
	db := models.Connection()
	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("%s", err)

		return
	}
	token := models.LoginCheck(db, input)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func GetUser(c *gin.Context) {
	if err := utils.TokenValid(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	id := c.Param("id")

	u := models.GetUser(id)

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})

}

func Delete(c *gin.Context) {
	if err := utils.TokenValid(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	id := c.Param("id")

	models.Delete(id)

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})

}

func Update(c *gin.Context) {
	if err := utils.TokenValid(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	id := c.Param("id")
	input := utils.RegisterType{}
	c.BindJSON(&input)
	u := models.Update(id, input)

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})

}

//crud operations for mongo

func RegisterUserMongo(c *gin.Context) {

	var input utils.RegisterType
	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("%s", err)

		return
	}
	fmt.Print(input)
	client := models.ConnectionMongo()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	input.Password = string(hashedPassword)

	models.SaveUser(input, client)
	c.JSON(http.StatusOK, gin.H{
		"message": "registered",
	})
}

func LoginMongo(c *gin.Context) {
	var input utils.RegisterType
	client := models.ConnectionMongo()
	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("%s", err)

		return
	}
	token := models.LoginCheckMongo(client, input)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func GetUserMongo(c *gin.Context) {
	if err := utils.TokenValid(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	username := c.Param("username")

	u := models.GetUserMongo(username)

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})

}

func UpdateMongo(c *gin.Context) {
	if err := utils.TokenValid(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	name := c.Param("name")
	input := utils.RegisterType{}
	c.BindJSON(&input)
	models.UpdateMongo(name, input)

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated",
	})

}

func DeleteMongo(c *gin.Context) {
	if err := utils.TokenValid(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	name := c.Param("name")

	models.DeleteMongo(name)

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})

}
