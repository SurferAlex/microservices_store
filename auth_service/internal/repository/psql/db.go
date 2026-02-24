package psql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
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

func RunMigrations(dbURL string) error {
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("Ошибка запуска миграций %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("База данных в актуальном состоянии")
		} else {
			return fmt.Errorf("критическая ошибка миграции: %w", err)
		}
	}
	return nil
}
