package main

import (
	"example/go_backoffice/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inizializzazione DB
	db := config.ConnectDB()

	// 2. Iniezione automatica tramite la funzione generata da Wire
	userController := InitUserController(db)

	// 3. Setup Router
	router := gin.Default()
	router.Use(config.CORS())

	router.GET("/users", userController.GetUsers)
	router.GET("/users/:id", userController.GetUserByID)
	router.POST("/users", userController.CreateUser)
	router.DELETE("/users/:id", userController.DeleteUser)

	router.Run("localhost:9090")
}
