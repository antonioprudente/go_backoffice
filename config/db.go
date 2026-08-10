package config

import (
	"example/go_backoffice/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql" // Cambia con gorm.io/driver/postgres se usi PostgreSQL
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	// Carica le variabili dal file .env se presente (es. sviluppo locale).
	// In produzione le variabili sono in genere già impostate dall'ambiente/orchestratore,
	// quindi ignoriamo l'errore se il file .env non esiste.
	if err := godotenv.Load(); err != nil {
		log.Println("Nessun file .env trovato, uso le variabili d'ambiente di sistema")
	}

	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "nome_db")

	// Crea il database se non esiste già, prima di connettersi ad esso.
	if err := ensureDatabaseExists(user, password, host, port, dbName); err != nil {
		log.Fatalf("Errore durante la creazione/verifica del database '%s': %v", dbName, err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Errore durante la connessione al DB: %v", err)
	}

	// Creazione/Aggiornamento automatico delle tabelle.
	// User va migrato insieme (o prima) agli altri due, perché AgentNode e
	// AgentOperator hanno FK verso User.
	err = db.AutoMigrate(
		&models.User{},
		&models.AgentNode{},
		&models.AgentOperator{},
	)
	if err != nil {
		log.Fatalf("Errore durante la migrazione del DB: %v", err)
	}

	return db
}

// ensureDatabaseExists si connette al server MySQL SENZA specificare un database
// (necessario perché connettersi direttamente a un DB inesistente fallirebbe),
// poi esegue CREATE DATABASE IF NOT EXISTS per crearlo automaticamente se manca.
func ensureDatabaseExists(user, password, host, port, dbName string) error {
	serverDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port,
	)

	serverDB, err := gorm.Open(mysql.Open(serverDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("impossibile connettersi al server MySQL: %w", err)
	}

	sqlDB, err := serverDB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	createStmt := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci",
		dbName,
	)
	if err := serverDB.Exec(createStmt).Error; err != nil {
		return fmt.Errorf("errore nella creazione del database: %w", err)
	}

	log.Printf("Database '%s' pronto (creato se non esisteva)", dbName)
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
