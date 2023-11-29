package main

import (
	"github.com/gin-gonic/gin"
	"github.com/raselldev/go-jwt/controllers"
	"github.com/raselldev/go-jwt/initializer"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectToDb()
	initializer.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.Run()
}
