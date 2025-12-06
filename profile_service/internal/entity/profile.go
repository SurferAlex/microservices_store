package entity

import "time"

type Profile struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id,omitempty" db:"user_id" binding:"omitempty,min=1"` // Изменить required на omitempty
	FirstName   string    `json:"first_name" db:"first_name" binding:"required,min=2,max=50"`
	LastName    string    `json:"last_name" db:"last_name" binding:"required,min=2,max=50"`
	Phone       string    `json:"phone" db:"phone" binding:"required,min=5,max=20"`
	AvatarURL   string    `json:"avatar_url,omitempty" db:"avatar_url" binding:"omitempty,url"`
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ProfileAddress struct {
	ID         int       `json:"id" db:"id"`
	ProfileID  int       `json:"profile_id,omitempty" db:"profile_id" binding:"omitempty,min=1"` // Убрать required
	Label      string    `json:"label,omitempty" db:"label" binding:"omitempty,max=50"`
	Country    string    `json:"country" db:"country" binding:"required,min=2,max=50"`
	City       string    `json:"city" db:"city" binding:"required,min=2,max=50"`
	Street     string    `json:"street" db:"street" binding:"required,min=2,max=100"`
	House      string    `json:"house" db:"house" binding:"required,min=1,max=20"`
	Apartment  string    `json:"apartment,omitempty" db:"apartment" binding:"omitempty,max=20"`
	PostalCode string    `json:"postal_code,omitempty" db:"postal_code" binding:"omitempty,max=20"`
	IsPrimary  bool      `json:"is_primary" db:"is_primary"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ProfileContact struct {
	ID         int       `json:"id" db:"id"`
	ProfileID  int       `json:"profile_id,omitempty" db:"profile_id" binding:"omitempty,min=1"` // Убрать required
	Type       string    `json:"type" db:"type" binding:"required,oneof=email phone messenger"`
	Value      string    `json:"value" db:"value" binding:"required,max=255"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
