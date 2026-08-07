package config

import (
	"example/go_backoffice/models"
	"log"

	"gorm.io/driver/mysql" // Cambia con gorm.io/driver/postgres se usi PostgreSQL
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := "root@tcp(127.0.0.1:3306)/nome_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Errore durante la connessione al DB: %v", err)
	}

	// Creazione/Aggiornamento automatico della tabella "users"
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Errore durante la migrazione del DB: %v", err)
	}

	return db
}
