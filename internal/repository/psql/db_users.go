package psql

import (
	"database/sql"
	"log"
	"vscode_test/internal/entity"
)

func GetUserByUsername(username string) (*entity.User, error) {
	log.Printf("Поиск пользователя: %s", username)
	var user entity.User

	query := `SELECT id, username, password FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Пользователь не найден: %s", username)
			return nil, nil
		}
		log.Printf("DB error: %v", err)
		return nil, err
	}
	log.Printf("Пользователь найден: %+v", user)
	return &user, nil
}

func InsertUser(user entity.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, user.Username, user.Email, user.Password)
	return err
}

func GetUserByID(id int) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, username, email, firstname, lastname, phone FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Phone)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func DeleteUser(username string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, username)
	return err
}
