package psql

import (
	"database/sql"
)

func AssignRoleToUser(userID int, roleName string) error {
	var roleID int
	// Получить id роли по имени
	err := db.QueryRow(`SELECT id FROM roles WHERE name=$1`, roleName).Scan(&roleID)
	if err != nil {
		return err
	}
	// Добавить связь роль-пользователь
	_, err = db.Exec(`INSERT INTO user_roles(user_id, role_id)
		VALUES ($1,$2) ON CONFLICT DO NOTHING`, userID, roleID)
	return err
}

func UserHasPermission(userID int, permission string) (bool, error) {
	query := `
	SELECT 1
	FROM user_roles ur
	JOIN role_permissions rp ON rp.role_id = ur.role_id
	JOIN permissions p ON p.id = rp.permission_id
	WHERE ur.user_id = $1 AND p.name = $2
	LIMIT 1`
	var one int
	err := db.QueryRow(query, userID, permission).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Получить слайс с названиями ролей пользователя
func GetUserRoles(userID int) ([]string, error) {
	rows, err := db.Query(`
		SELECT r.name
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		roles = append(roles, name)
	}
	return roles, nil
}
