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
	agencyOperatorController := InitAgencyOperatorController(db)
	agentOperatorController := InitAgentOperatorController(db)

	// Setup Router
	router := gin.Default()
	router.Use(config.CORS())

	router.POST("/login", authController.Login)

	protected := router.Group("")

	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/logout", authController.Logout) //POST /logout

		// OPERATORS + PIVOT CALLS
		operators := protected.Group("/operators")
		pivot := protected.Group("/operator")
		operators.Use(
			middlewares.RequireRoles(enums.RoleAdmin.String()),
			middlewares.SetRoleMiddleware(enums.RoleOperator.String()),
		)
		{
			operators.POST("", userController.CreateUser) // POST /operators
			operators.GET("", userController.GetUsers)    // GET /operators

			pivot.POST("/to/agency", agencyOperatorController.AssignAgencyToOperator) // POST /operator/to/agency
			pivot.POST("/to/agent", agentOperatorController.AssignAgentToOperator)    // POST /operator/to/agent
			pivot.POST("/to/agencies")                                                // POST /operator/to/agencies
			pivot.POST("/to/agents")                                                  // POST /operator/to/agents
		}

		// AGENTS CALLS
		agents := protected.Group("/agents")
		agents.Use(
			middlewares.RequireRoles(enums.RoleAdmin.String(), enums.RoleOperator.String()),
			middlewares.SetRoleMiddleware(enums.RoleAgent.String()),
		)
		{
			operators.GET("/:id", userController.GetUserByID)  // GET /operators/:id
			agents.POST("", agentController.CreateAgentNode)   // POST /agents
			agents.GET("", userController.GetUsers)            // GET /agents
			agents.GET("/:id", userController.GetUserByID)     // GET /agents/:id
			agents.GET("/tree", agentController.GetAgentsTree) // GET /agents/tree
		}

		// AGENCIES CALLS
		agencies := protected.Group("/agencies")
		agencies.Use(
			middlewares.RequireRoles(enums.RoleAdmin.String(), enums.RoleOperator.String()),
			middlewares.SetRoleMiddleware(enums.RoleAgency.String()),
		)
		{
			agencies.POST("", userController.CreateUser)     // POST /agencies
			agencies.GET("", userController.GetUsers)        // GET /agencies
			agencies.GET("/:id", userController.GetUserByID) // GET /agencies/:id
		}

		// USERS CALLS
		users := protected.Group("/users")
		users.Use(
			middlewares.RequireRoles(enums.RoleAdmin.String(), enums.RoleOperator.String()),
			middlewares.SetRoleMiddleware(enums.RoleUser.String()),
		)
		{
			users.POST("", userController.CreateUser)                   // POST /users
			users.GET("", userController.GetUsers)                      // GET /users
			users.GET("/:id", userController.GetUserByID)               // GET /users/:id
			users.DELETE("/:id", userController.DeleteUser)             // DELETE /users/:id
			users.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
			users.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
			users.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
		}
	}

	router.Run("localhost:9090")
}
