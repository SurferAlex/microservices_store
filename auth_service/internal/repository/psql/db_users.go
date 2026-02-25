package psql

import (
	"auth_service/internal/entity"
	"database/sql"
	"log"
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

func GetUserByEmail(email string) (*entity.User, error) {
	row := db.QueryRow(`SELECT id, username, email, password FROM users WHERE email = $1`, email)
	var u entity.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}

func InsertUser(user entity.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, user.Username, user.Email, user.Password)
	return err
}

func GetUserByID(id int) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, username, email FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
