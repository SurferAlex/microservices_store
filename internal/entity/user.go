package entity

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	Address   Address   `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Address struct {
	Country    string `json:"country"`
	City       string `json:"city"`
	Street     string `json:"street"`
	House      string `json:"house"`
	Apartment  string `json:"apartment"`
	PostalCode string `json:"postal_code"`
}
