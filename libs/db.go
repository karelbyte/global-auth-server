package libs

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

func GetDB() *sql.DB {
	once.Do(func() {
		_ = godotenv.Load()

		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		database := os.Getenv("DB_NAME")

		connString := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s", host, port, user, password, database)

		var err error
		dbInstance, err = sql.Open("sqlserver", connString)
		if err != nil {
			panic(fmt.Sprintf("Error connecting to the database: %s", err))
		}
		if err = dbInstance.Ping(); err != nil {
			panic(fmt.Sprintf("Ping could not be done to the database: %s", err))
		}
	})
	return dbInstance
}
