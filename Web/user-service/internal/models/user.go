package models

import "time"

type User struct {
	ID           int
	Login        string
	FirstName    string
	LastName     string
	BirthDate    time.Time
	PasswordHash string
	CreatedAt    time.Time
}
