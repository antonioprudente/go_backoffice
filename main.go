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

	// Crea l'utente ADMIN di default se non esiste già
	config.SeedAdminUser(db)
	// Iniezione controllers automatica tramite la funzione generata da Wire
	userController := InitUserController(db)
	agentController := InitAgentController(db)
	authController := InitAuthController(db)

	// Setup Router
	router := gin.Default()
	router.Use(config.CORS())

	router.POST("/login", authController.Login)

	protected := router.Group("/api")
	protected.Use(middlewares.AuthRequired())
	{
		users := protected.Group("/users")
		users.Use(middlewares.SetRoleMiddleware(enums.RoleUser.String()))
		{
			users.POST("", userController.CreateUser)                   // POST /api/users
			users.GET("", userController.GetUsers)                      // GET /api/users
			users.GET("/:id", userController.GetUserByID)               // GET /api/users/:id
			users.DELETE("/:id", userController.DeleteUser)             // DELETE /api/users/:id
			users.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /api/users/:id/active
			users.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /api/users/:id/suspend
			users.PATCH("/:id/block", userController.BlockUserById)     // PATCH /api/users/:id/block
		}

		agents := protected.Group("/agents")
		agents.Use(middlewares.SetRoleMiddleware(enums.RoleAgent.String()))
		{
			agents.POST("", agentController.CreateAgentNode)   // POST /api/agents
			agents.GET("", userController.GetUsers)            // GET /api/agents
			agents.GET("/:id", userController.GetUserByID)     // GET /api/agents/:id
			agents.GET("/tree", agentController.GetAgentsTree) // GET /api/agents/tree
		}

		agencies := protected.Group("/agencies")
		agencies.Use(middlewares.SetRoleMiddleware(enums.RoleAgency.String()))
		{
			agencies.POST("", userController.CreateUser)     // POST /api/agencies
			agencies.GET("", userController.GetUsers)        // GET /api/agencies
			agencies.GET("/:id", userController.GetUserByID) // GET /api/agencies/:id
		}
	}

	router.Run("localhost:9090")
}
