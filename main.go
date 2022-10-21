package main

import (
	_ "bytes"
	"crud-go/controllers"
	_ "database/sql"
	_ "fmt"
	_ "log"

	"github.com/gin-gonic/gin"
)

// var body2 bytes.Buffer

func main() {

	r := gin.Default()
	public := r.Group("/api")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	public.GET("/getUser/:id", controllers.GetUser)
	public.DELETE("/delete/:id", controllers.Delete)
	public.PATCH("/update/:id", controllers.Update)
	public.POST("/insertElastic", controllers.InsertElastic)
	public.GET("/getElastic", controllers.Get)
	public.POST("/registerMongo", controllers.RegisterUserMongo)
	public.POST("/loginMongo", controllers.LoginMongo)
	public.GET("/getUserMongo/:username", controllers.GetUserMongo)
	public.DELETE("/deleteMongo/:id", controllers.DeleteMongo)
	public.PATCH("/updateMongo/:name", controllers.UpdateMongo)
	r.Run(":8000")
}
