package libs

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
		_ = godotenv.Load()

		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		database := os.Getenv("DB_NAME")

		// Variación 1: Sin SSL
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, database)

		// Variación 2: Con SSL y parámetros adicionales
		connString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require connect_timeout=10",
			host, port, user, password, database)

		// Variación 3: Con SSL y verificación de certificado desactivada
		connString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require sslrootcert=/path/to/rds-ca-2019-root.pem",
			host, port, user, password, database)

		fmt.Printf("Intentando conectar a PostgreSQL en %s:%s\n", host, port)

		var err error
		dbInstance, err = sql.Open("postgres", connString)
		if err != nil {
			panic(fmt.Sprintf("Error connecting to database: %s", err))
		}

		// Configurar el pool de conexiones
		dbInstance.SetMaxOpenConns(25)
		dbInstance.SetMaxIdleConns(5)
		dbInstance.SetConnMaxLifetime(5 * time.Minute)

		// Verificar la conexión
		if err = dbInstance.Ping(); err != nil {
			fmt.Printf("Error al hacer ping: %v\n", err)
			panic(fmt.Sprintf("Could not ping database: %s", err))
		}
		fmt.Println("Conexión exitosa a PostgreSQL")
	})
	return dbInstance
}
