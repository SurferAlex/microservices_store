package psql

import (
	"database/sql"
	"fmt"
	"profile_service/internal/entity"
)

const addressesTable = "profile_addresses"

func CreateAddress(address entity.ProfileAddress) (*entity.ProfileAddress, error) {
	// Если устанавливается основной адрес, сначала сбрасываем все остальные
	if address.IsPrimary {
		queryReset := `
			UPDATE ` + addressesTable + `
			SET is_primary = false
			WHERE profile_id = $1;
		`
		_, err := db.Exec(queryReset, address.ProfileID)
		if err != nil {
			return nil, err
		}
	}

	query := `
		INSERT INTO ` + addressesTable + ` (profile_id, label, country, city, street, house, apartment, postal_code, is_primary)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at;
	`

	err := db.QueryRow(
		query,
		address.ProfileID,
		address.Label,
		address.Country,
		address.City,
		address.Street,
		address.House,
		address.Apartment,
		address.PostalCode,
		address.IsPrimary,
	).Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &address, nil
}

func GetAddressByID(id int) (*entity.ProfileAddress, error) {
	query := `
		SELECT id, profile_id, label, country, city, street, house, apartment, postal_code, is_primary, created_at, updated_at
		FROM ` + addressesTable + `
		WHERE id = $1;
	`

	var address entity.ProfileAddress
	err := db.QueryRow(query, id).Scan(
		&address.ID,
		&address.ProfileID,
		&address.Label,
		&address.Country,
		&address.City,
		&address.Street,
		&address.House,
		&address.Apartment,
		&address.PostalCode,
		&address.IsPrimary,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &address, nil
}

func GetAddressesByProfileID(profileID int) ([]entity.ProfileAddress, error) {
	query := `
	     SELECT id, profile_id, label, country, city, street, house, apartment, postal_code, is_primary, created_at, updated_at
         FROM profile_addresses
         WHERE profile_id = $1
         ORDER BY created_at;
	 `

	rows, err := db.Query(query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []entity.ProfileAddress

	for rows.Next() {
		var addr entity.ProfileAddress
		err := rows.Scan(
			&addr.ID,
			&addr.ProfileID,
			&addr.Label,
			&addr.Country,
			&addr.City,
			&addr.Street,
			&addr.House,
			&addr.Apartment,
			&addr.PostalCode,
			&addr.IsPrimary,
			&addr.CreatedAt,
			&addr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return addresses, nil

}

func UpdateAddress(address entity.ProfileAddress) error {
	// Если устанавливается основной адрес, сначала сбрасываем все остальные
	if address.IsPrimary {
		queryReset := `
			UPDATE ` + addressesTable + `
			SET is_primary = false
			WHERE profile_id = $1 AND id != $2;
		`
		_, err := db.Exec(queryReset, address.ProfileID, address.ID)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE ` + addressesTable + `
		SET label = $1,
		    country = $2,
		    city = $3,
		    street = $4,
		    house = $5,
		    apartment = $6,
		    postal_code = $7,
		    is_primary = $8,
		    updated_at = NOW()
		WHERE id = $9;
	`

	result, err := db.Exec(
		query,
		address.Label,
		address.Country,
		address.City,
		address.Street,
		address.House,
		address.Apartment,
		address.PostalCode,
		address.IsPrimary,
		address.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("адрес с id %d не найден", address.ID)
	}

	return nil
}

func DeleteAddress(id int) error {
	query := `
		DELETE FROM ` + addressesTable + ` WHERE id = $1;
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
		return fmt.Errorf("адрес с id %d не найден", id)
	}

	return nil
}

func SetPrimaryAddress(profileID int, addressID int) error {
	// Проверяем существование профиля
	_, err := GetProfileByID(profileID)
	if err != nil {
		return fmt.Errorf("профиль с id %d не найден", profileID)
	}

	// Проверяем существование адреса и что он принадлежит профилю
	address, err := GetAddressByID(addressID)
	if err != nil {
		return fmt.Errorf("ошибка при получении адреса: %w", err)
	}
	if address == nil {
		return fmt.Errorf("адрес с id %d не найден", addressID)
	}
	if address.ProfileID != profileID {
		return fmt.Errorf("адрес с id %d не принадлежит профилю %d", addressID, profileID)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Сбрасываем все адреса профиля
	query1 := `
		UPDATE ` + addressesTable + `
		SET is_primary = false
		WHERE profile_id = $1;
	`
	_, err = tx.Exec(query1, profileID)
	if err != nil {
		return err
	}

	// Устанавливаем основной адрес
	query2 := `
		UPDATE ` + addressesTable + `
		SET is_primary = true, updated_at = NOW()
		WHERE id = $1 AND profile_id = $2;
	`
	result, err := tx.Exec(query2, addressID, profileID)
	if err != nil {
		return err
	}

	// Проверяем, что адрес был обновлен
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("не удалось установить основной адрес")
	}

	return tx.Commit()
}
