package main

import (
	"example/go_backoffice/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inizializzazione DB
	db := config.ConnectDB()

	// Iniezione controllers automatica tramite la funzione generata da Wire
	userController := InitUserController(db)

	// Setup Router
	router := gin.Default()
	router.Use(config.CORS())

	users := router.Group("/users")
	{
		users.GET("", userController.GetUsers)    // GET /users
		users.POST("", userController.CreateUser) // POST /users

		users.GET("/:id", userController.GetUserByID)   // GET /users/:id
		users.DELETE("/:id", userController.DeleteUser) // DELETE /users/:id

		users.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
		users.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
		users.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
	}

	router.Run("localhost:9090")
}
