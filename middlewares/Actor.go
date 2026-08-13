package middlewares

import (
	"errors"

	"example/go_backoffice/policies"

	"github.com/gin-gonic/gin"
)

// ActorFromContext costruisce l'AuthContext leggendo userID/userRole
// impostati da AuthMiddleware. Va chiamato solo su rotte protette da AuthMiddleware.
func ActorFromContext(ctx *gin.Context) (policies.AuthContext, error) {
	idVal, exists := ctx.Get("userID")
	if !exists {
		return policies.AuthContext{}, errors.New("utente non autenticato")
	}
	roleVal, exists := ctx.Get("userRole")
	if !exists {
		return policies.AuthContext{}, errors.New("utente non autenticato")
	}

	userID, ok := idVal.(uint)
	if !ok {
		return policies.AuthContext{}, errors.New("userID non valido nel context")
	}
	role, ok := roleVal.(string)
	if !ok {
		return policies.AuthContext{}, errors.New("userRole non valido nel context")
	}

	return policies.AuthContext{UserID: userID, Role: role}, nil
}
