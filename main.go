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

		//protected.POST("/logout", authController.Logout) //POST /logout (non ancora implementato)

		// OPERATORS
		operators := protected.Group("/operators")
		{
			// Rotte per ADMIN
			adminOnly := operators.Group("")
			adminOnly.Use(
				middlewares.RequireRoles(enums.RoleAdmin.String()),
				middlewares.SetRoleMiddleware(enums.RoleOperator.String()),
			)
			{
				adminOnly.POST("", userController.CreateUser)
				adminOnly.GET("", userController.GetUsers)
				adminOnly.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
				adminOnly.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
				adminOnly.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
			}

			// Rotta per ADMIN e OPERATOR (self-check/policy)
			selfAccess := operators.Group("")
			selfAccess.Use(
				middlewares.RequireRoles(enums.RoleAdmin.String(), enums.RoleOperator.String()),
				middlewares.SetRoleMiddleware(enums.RoleOperator.String()),
			)
			{
				selfAccess.GET("/:id", userController.GetUserByID)
			}
		}

		// PIVOT
		pivot := protected.Group("/operator")
		pivot.Use(
			middlewares.RequireRoles(enums.RoleAdmin.String()),
		)
		{
			pivot.POST("/to/agency", agencyOperatorController.AssignAgencyToOperator)
			pivot.POST("/to/agent", agentOperatorController.AssignAgentToOperator)
			pivot.POST("/to/agencies", agencyOperatorController.AssignAgenciesToOperator)
			pivot.POST("/to/agents", agentOperatorController.AssignAgentsToOperator)
			pivot.DELETE("/:operatorID/from/agent/:agentID", agentOperatorController.RemoveAgentFromOperator)
		}

		// AGENTS CALLS
		agents := protected.Group("/agents")
		agents.Use(
			middlewares.RequireRoles(enums.RoleOperator.String(), enums.RoleAdmin.String(), enums.RoleAgent.String()),
			middlewares.SetRoleMiddleware(enums.RoleAgent.String()),
		)
		{
			agents.POST("", agentController.CreateAgentNode)             // POST /agents
			agents.GET("", userController.GetUsers)                      // GET /agents
			agents.GET("/:id", userController.GetUserByID)               // GET /agents/:id
			agents.GET("/tree", agentController.GetFilteredTree)         // GET /agents/tree
			agents.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
			agents.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
			agents.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
		}

		// AGENCIES CALLS
		agencies := protected.Group("/agencies")
		agencies.Use(
			middlewares.RequireRoles(enums.RoleAdmin.String(), enums.RoleOperator.String(), enums.RoleAgent.String()),
			middlewares.SetRoleMiddleware(enums.RoleAgency.String()),
		)
		{
			agencies.POST("", userController.CreateUser)                   // POST /agencies
			agencies.GET("", userController.GetUsers)                      // GET /agencies
			agencies.GET("/:id", userController.GetUserByID)               // GET /agencies/:id
			agencies.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
			agencies.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
			agencies.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
		}

		// USERS CALLS
		users := protected.Group("/users")
		users.Use(
			middlewares.RequireRoles(enums.RoleAdmin.String(), enums.RoleOperator.String(), enums.RoleAgent.String(), enums.RoleAgency.String()),
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
