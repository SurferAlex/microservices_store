package psql

import (
	"database/sql"
	"fmt"
	"profile_service/internal/entity"
)

const profilesTable = "profiles"

func CreateProfile(profile entity.Profile) (*entity.Profile, error) {
	query := `
		INSERT INTO ` + profilesTable + ` (user_id, first_name, last_name, phone, avatar_url, date_of_birth)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at;
	`

	err := db.QueryRow(
		query,
		profile.UserID,
		profile.FirstName,
		profile.LastName,
		profile.Phone,
		profile.AvatarURL,
		profile.DateOfBirth,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func GetProfileByUserID(userID int) (*entity.Profile, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone, avatar_url, date_of_birth, created_at, updated_at
		FROM ` + profilesTable + `
		WHERE user_id = $1;
	`

	var profile entity.Profile
	err := db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
		&profile.AvatarURL,
		&profile.DateOfBirth,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func GetProfileByID(id int) (*entity.Profile, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone, avatar_url, date_of_birth, created_at, updated_at
		FROM ` + profilesTable + `
		WHERE id = $1;
	`

	var profile entity.Profile
	err := db.QueryRow(query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
		&profile.AvatarURL,
		&profile.DateOfBirth,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Профиль не найден
		}
		return nil, err // Другая ошибка БД
	}

	return &profile, nil
}

func UpdateProfile(profile entity.Profile) error {
	query := `
		UPDATE ` + profilesTable + `
		SET first_name = $1,
		    last_name = $2,
		    phone = $3,
		    avatar_url = $4,
		    date_of_birth = $5,
		    updated_at = NOW()
		WHERE id = $6;
	`

	_, err := db.Exec(
		query,
		profile.FirstName,
		profile.LastName,
		profile.Phone,
		profile.AvatarURL,
		profile.DateOfBirth,
		profile.ID,
	)
	return err
}

func DeleteProfile(id int) error {
	query := `
		DELETE FROM profiles WHERE id = $1;
	`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("профиль с id %d не найден", id)
	}

	return nil
}
