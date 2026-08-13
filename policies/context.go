package policies

// AuthContext rappresenta l'utente autenticato che sta eseguendo l'azione,
// estratto dal JWT. Non dipende da gin: resta un package di dominio puro.
type AuthContext struct {
	UserID uint
	Role   string
}
