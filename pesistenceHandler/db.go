package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var psqlInfo = fmt.Sprintf("host=%s port=%s user=%s "+
	"dbname=%s sslmode=disable", GetEnvVariable("PGHOST"), GetEnvVariable("PGPORT"), GetEnvVariable("PGUSER"), GetEnvVariable("PGDBNAME"))

var db *sql.DB

// InitPostgresDB function returns the DB connection pointer
func InitPostgresDB() *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to Postgres DB successfully!!")
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		panic(err)
	}
	return db
}
