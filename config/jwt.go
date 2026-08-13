package config

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims personalizzate: ci mettiamo l'ID utente e il ruolo
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	secret := GetEnv("JWT_SECRET", "dev-secret-cambiami") // in produzione DEVE arrivare da env
	return []byte(secret)
}

// GenerateToken crea un JWT firmato per l'utente loggato
func GenerateToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

// ParseToken valida il token e ritorna le claims
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metodo di firma inatteso")
		}
		return jwtSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token non valido")
	}
	return claims, nil
}
