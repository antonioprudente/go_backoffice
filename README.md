readme_content = """# 🚀 Go Backoffice REST API

Un'applicazione backend RESTful sviluppata in **Go (Golang)** utilizzando il framework **Gin**. Il progetto segue un'architettura layered (a livelli) per garantire scalabilità, manutenibilità e una netta separazione delle responsabilità (Separation of Concerns).

---

## 🛠️ Tech Stack & Dipendenze

- **Linguaggio**: Go (Golang)
- **Web Framework**: Gin Web Framework
- **Database**: SQL (Database relazionale con driver `database/sql`)
- **Autenticazione**: JSON Web Token (`golang-jwt/jwt`)
- **Iniezione Dipendenze**: Google Wire
- **Middleware**: CORS & Role-Based Access Control (RBAC)

---

## 📁 Struttura del Progetto

Il progetto è organizzato secondo il pattern **Layered Architecture**:

```text
go_backoffice/
├── config/             # Configurazione del database, CORS e variabili d'ambiente
│   ├── db.go           # Connessione al DB SQL
│   └── cors.go         # Gestione delle policy CORS
├── controllers/        # Gestione delle richieste HTTP (Gin Handler)
│   ├── user_controller.go
│   └── agent_controller.go
├── services/           # Logica di business e validazione dati
│   ├── user_service.go
│   └── agent_service.go
├── repositories/       # Query SQL e interazione diretta con il database
│   ├── user_repository.go
│   └── agent_repository.go
├── models/             # Strutture dati (Struct)
│   ├── user.go
│   └── agent.go
├── middlewares/        # Middleware per Auth, CORS e Gestione Ruoli
│   ├── AuthMiddleware.go
│   └── RoleMiddleware.go
├── routes/             # Definizioni dei gruppi di rotte e applicazione middleware
│   └── router.go
├── main.go             # Entrypoint dell'applicazione
├── go.mod
└── .env                # Variabili d'ambiente (non committato)