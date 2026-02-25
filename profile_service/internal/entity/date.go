package entity

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Date - кастомный тип для даты без времени
type Date struct {
	time.Time
}

// UnmarshalJSON парсит дату из JSON формата "2006-01-02"
func (d *Date) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" || str == "null" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %w", err)
	}
	d.Time = t
	return nil
}

// MarshalJSON сериализует дату в JSON формат "2006-01-02"
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Time.Format("2006-01-02"))), nil
}

// Value для работы с БД
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// Scan для работы с БД
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		return d.UnmarshalJSON(v)
	case string:
		return d.UnmarshalJSON([]byte(v))
	default:
		return fmt.Errorf("неверный тип для Date: %T", value)
	}
}
