package psql

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения к БД
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Успешное подключение к БД", dataSourceName)
}

func GetDB() *sql.DB {
	return db
}
