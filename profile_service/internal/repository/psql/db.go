package psql

import (
	"database/sql"
	"fmt"
	"log"
)

var (
	db *sql.DB
)

func InitDB(connStr string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть соединение: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения БД: %w", err)
	}
	// Логирование подключения
	log.Printf("✅ Успешное подключение к БД: %v", db)

	// Проверка текущей БД
	var currentDB string
	err = db.QueryRow("SELECT current_database()").Scan(&currentDB)
	if err == nil {
		log.Printf("📊 Текущая БД: %s", currentDB)
	}

	return db, nil

}
