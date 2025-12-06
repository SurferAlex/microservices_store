package psql

import (
	"fmt"
	"profile_service/internal/entity"
)

const contactsTable = "profile_contacts"

func CreateContact(contact entity.ProfileContact) (*entity.ProfileContact, error) {
	query := `
		INSERT INTO ` + contactsTable + ` (profile_id, type, value, is_verified)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at;
	`

	err := db.QueryRow(
		query,
		contact.ProfileID,
		contact.Type,
		contact.Value,
		contact.IsVerified,
	).Scan(&contact.ID, &contact.CreatedAt, &contact.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func GetContactsByProfileID(profileID int) ([]entity.ProfileContact, error) {
	query := `
		SELECT id, profile_id, type, value, is_verified, created_at, updated_at
		FROM ` + contactsTable + `
		WHERE profile_id = $1
		ORDER BY created_at;
	`

	rows, err := db.Query(query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []entity.ProfileContact

	for rows.Next() {
		var contact entity.ProfileContact
		err := rows.Scan(
			&contact.ID,
			&contact.ProfileID,
			&contact.Type,
			&contact.Value,
			&contact.IsVerified,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func UpdateContact(contact entity.ProfileContact) error {
	query := `
		UPDATE ` + contactsTable + `
		SET type = $1,
		    value = $2,
		    is_verified = $3,
		    updated_at = NOW()
		WHERE id = $4;
	`

	result, err := db.Exec(
		query,
		contact.Type,
		contact.Value,
		contact.IsVerified,
		contact.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("контакт с id %d не найден", contact.ID)
	}

	return nil
}

func DeleteContact(id int) error {
	query := `
		DELETE FROM ` + contactsTable + ` WHERE id = $1;
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
		return fmt.Errorf("контакт с id %d не найден", id)
	}

	return nil
}

// Опционально: подтвердить контакт
func VerifyContact(id int) error {
	query := `
		UPDATE ` + contactsTable + `
		SET is_verified = true,
		    updated_at = NOW()
		WHERE id = $1;
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
		return fmt.Errorf("контакт с id %d не найден", id)
	}

	return nil
}
