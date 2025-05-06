package libs

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var (
	dbInstance *sql.DB
	once       sync.Once
)

// Carga la configuración desde variables de entorno
func LoadDBConfigFromEnv() DBConfig {
	_ = godotenv.Load()

	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

// Set a new connection to the database using the conference
func NewDB(cfg DBConfig) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"server=%s;port=%s;user id=%s;password=%s;database=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error abriendo la conexión: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("no se pudo hacer ping a la base de datos: %w", err)
	}

	return db, nil
}

// Singleton: devuelve una instancia compartida
func GetDB() *sql.DB {
	once.Do(func() {
		cfg := LoadDBConfigFromEnv()
		db, err := NewDB(cfg)
		if err != nil {
			panic(err) // sigue siendo crítico si falla la conexión
		}
		dbInstance = db
	})
	return dbInstance
}
