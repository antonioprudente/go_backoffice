package config

import (
	"example/go_backoffice/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPGDB() *gorm.DB {
	// Carica le variabili dal file .env se presente
	if err := godotenv.Load(); err != nil {
		log.Println("Nessun file .env trovato, uso le variabili d'ambiente di sistema")
	}

	user := os.Getenv("PG_DB_USER")
	password := os.Getenv("PG_DB_PASSWORD")
	host := os.Getenv("PG_DB_HOST")
	port := os.Getenv("PG_DB_PORT")
	dbName := os.Getenv("PG_DB_NAME")

	// 1. Assicura la presenza del database usando una connessione al DB di default 'postgres'
	if err := ensureDatabasePGExists(user, password, host, port, dbName); err != nil {
		log.Fatalf("Errore durante la creazione/verifica del database '%s': %v", dbName, err)
	}

	// 2. DSN corretto per driver postgres
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Rome",
		host, user, password, dbName, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		log.Fatalf("Errore durante la connessione al DB: %v", err)
	}

	db.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN CREATE TYPE role AS ENUM ('ADMIN', 'OPERATOR', 'AGENT', 'AGENCY', 'USER'); END IF; END $$;")
	db.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN CREATE TYPE status AS ENUM ('ACTIVE', 'SUSPENDED', 'BLOCKED', 'DEFAULT'); END IF; END $$;")
	db.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'action') THEN CREATE TYPE action AS ENUM ('CREATE', 'UPDATE', 'DELETE', 'RESTORE', 'ASSIGNMENT', 'REMOVE', 'MOVE', 'ACTIVE', 'SUSPEND', 'BLOCK'); END IF; END $$;")

	// 3. Creazione/Aggiornamento automatico delle tabelle
	err = db.AutoMigrate(
		&models.User{},
		&models.AgentNode{},
		&models.AgentOperator{},
		&models.AgencyOperator{},
		&models.Note{},
		&models.ActivityLog{},
	)
	if err != nil {
		log.Fatalf("Errore durante la migrazione del DB: %v", err)
	}

	return db
}

// ensureDatabasePGExists si connette al DB di sistema 'postgres' per verificare
// se il database target esiste ed eventualmente crearlo.
func ensureDatabasePGExists(user, password, host, port, dbName string) error {
	// Connessione al database amministrativo 'postgres'
	systemDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		host, user, password, port,
	)

	systemDB, err := gorm.Open(postgres.Open(systemDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("impossibile connettersi al server PostgreSQL: %w", err)
	}

	sqlDB, err := systemDB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// Verifica se il database esiste già
	var count int
	checkQuery := "SELECT count(*) FROM pg_database WHERE datname = ?"
	if err := systemDB.Raw(checkQuery, dbName).Scan(&count).Error; err != nil {
		return fmt.Errorf("errore durante la verifica dell'esistenza del DB: %w", err)
	}

	// Se il DB non esiste, eseguibile CREATE DATABASE
	if count == 0 {
		createStmt := fmt.Sprintf("CREATE DATABASE %s", dbName)
		if err := systemDB.Exec(createStmt).Error; err != nil {
			return fmt.Errorf("errore nella creazione del database: %w", err)
		}
		log.Printf("Database '%s' creato con successo", dbName)
	} else {
		log.Printf("Database '%s' già esistente", dbName)
	}

	return nil
}
