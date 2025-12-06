package psql

import (
	"database/sql"
	"log"
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

	// Логирование подключения
	log.Printf("✅ Успешное подключение к БД: %s", dataSourceName)

	// Проверка текущей БД
	var currentDB string
	err = db.QueryRow("SELECT current_database()").Scan(&currentDB)
	if err == nil {
		log.Printf("📊 Текущая БД: %s", currentDB)
	}
}
