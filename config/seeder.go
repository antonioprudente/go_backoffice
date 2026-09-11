package config

import (
	"log"
	"os"

	"example/go_backoffice/enums"
	"example/go_backoffice/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdminUser crea automaticamente un utente con ruolo ADMIN al primo avvio,
// se non ne esiste già uno nel database. Le credenziali sono configurabili
// tramite variabili d'ambiente (utile per non hardcodare la password).
func SeedAdminUser(db *gorm.DB) {
	var count int64
	if err := db.Model(&models.User{}).Where("role = ?", enums.RoleAdmin).Count(&count).Error; err != nil {
		log.Printf("Errore durante la verifica dell'utente ADMIN: %v", err)
		return
	}

	if count > 0 {
		// Esiste già almeno un ADMIN, non serve fare nulla
		return
	}

	email := os.Getenv("ADMIN_EMAIL")
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	firstName := os.Getenv("ADMIN_FIRST_NAME")
	lastName := os.Getenv("ADMIN_LAST_NAME")

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Errore durante l'hashing della password ADMIN: %v", err)
		return
	}

	admin := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Role:      enums.RoleAdmin,
		Status:    enums.StatusActive,
		Email:     email,
		Password:  string(hashed),
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Errore durante la creazione dell'utente ADMIN: %v", err)
		return
	}

	log.Printf("Utente ADMIN creato con successo (email: %s, username: %s)", email, username)
}
