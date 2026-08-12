package main

import (
	"example/go_backoffice/config"
	"example/go_backoffice/enums"
	"example/go_backoffice/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inizializzazione DB
	db := config.ConnectDB()

	// Iniezione controllers automatica tramite la funzione generata da Wire
	userController := InitUserController(db)
	agentController := InitAgentController(db)

	// Setup Router
	router := gin.Default()
	router.Use(config.CORS())

	users := router.Group("/users")
	{
		users.Use(middlewares.SetRoleMiddleware(enums.RoleUser.String()))
		users.GET("", userController.GetUsers)                      // GET /users
		users.POST("", userController.CreateUser)                   // POST /users
		users.GET("/:id", userController.GetUserByID)               // GET /users/:id
		users.DELETE("/:id", userController.DeleteUser)             // DELETE /users/:id
		users.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
		users.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
		users.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
	}

	agents := router.Group("/agents")
	{
		agents.Use(middlewares.SetRoleMiddleware(enums.RoleAgent.String()))
		agents.GET("/tree", agentController.GetAgentsTree) // GET /agents/tree
		agents.POST("", agentController.CreateAgentNode)   // POST /agents
		agents.GET("", userController.GetUsers)            // GET /agents
	}
	agencies := router.Group("/agencies")
	{
		agencies.Use(middlewares.SetRoleMiddleware(enums.RoleAgency.String()))
		agencies.POST("", userController.CreateUser)
		agencies.GET("", userController.GetUsers)
		agencies.GET("/:id", userController.GetUserByID)
	}

	router.Run("localhost:9090")
}
