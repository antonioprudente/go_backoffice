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

		protected.POST("/logout", authController.Logout) //POST /logout (non ancora implementato)

		protected.GET("/profile", userController.GetPersonalProfile)

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
				adminOnly.POST("", userController.CreateUser)                   // POST /operators
				adminOnly.GET("", userController.GetUsers)                      // GET /operators
				adminOnly.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
				adminOnly.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
				adminOnly.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
				adminOnly.DELETE("/:id", userController.DeleteUser)             // DELETE /users/:id
			}

			// Rotta per ADMIN e OPERATOR (self-check/policy)
			selfAccess := operators.Group("")
			selfAccess.Use(
				middlewares.RequireRoles(enums.RoleAdmin.String(), enums.RoleOperator.String()),
				middlewares.SetRoleMiddleware(enums.RoleOperator.String()),
			)
			{
				selfAccess.GET("/:id", userController.GetUserByID) // GET /operators/:id
				selfAccess.PUT("/:id", userController.UpdateUser)  // PUT /operators
			}
		}

		// PIVOT
		pivot := protected.Group("/operator")
		pivot.Use(
			middlewares.RequireRoles(enums.RoleAdmin.String()),
		)
		{
			pivot.POST("/to/agency", agencyOperatorController.AssignAgencyToOperator)                             // POST /operator/to/agency
			pivot.POST("/to/agent", agentOperatorController.AssignAgentToOperator)                                // POST /operator/to/agent
			pivot.POST("/to/agencies", agencyOperatorController.AssignAgenciesToOperator)                         // POST /operator/to/agencies
			pivot.POST("/to/agents", agentOperatorController.AssignAgentsToOperator)                              // POST /operator/to/agents
			pivot.DELETE("/:operatorID/from/agency/:agencyID", agencyOperatorController.RemoveAgencyFromOperator) // DELETE /operator/:operatorID/from/agency/:agencyID
			pivot.DELETE("/:operatorID/from/agent/:agentID", agentOperatorController.RemoveAgentFromOperator)     // DELETE /operator/:operatorID/from/agent/:agentID
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
			agents.PUT("/:id", userController.UpdateUser)                // PUT /agents
			agents.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
			agents.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
			agents.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
			agents.DELETE("/:id", agentController.DeleteAgent)           // DELETE /agents/:id
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
			agencies.PUT("/:id", userController.UpdateUser)                // PUT /agencies
			agencies.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
			agencies.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
			agencies.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
			agencies.DELETE("/:id", userController.DeleteUser)             // DELETE /agencies/:id
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
			users.PUT("/:id", userController.UpdateUser)                // PUT /users
			users.PATCH("/:id/active", userController.ActiveUserById)   // PATCH /users/:id/active
			users.PATCH("/:id/suspend", userController.SuspendUserById) // PATCH /users/:id/suspend
			users.PATCH("/:id/block", userController.BlockUserById)     // PATCH /users/:id/block
			users.DELETE("/:id", userController.DeleteUser)             // DELETE /users/:id
		}
	}

	router.Run("localhost:9090")
}
