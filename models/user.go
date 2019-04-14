package models

import (
	"database/sql"
	"goa/utils/crypto"
)

type User struct {
	// Email of the user
	Email *string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
	// First name of the user
	FirstName *string `form:"first_name,omitempty" json:"first_name,omitempty" xml:"first_name,omitempty"`
	// ID of account
	ID *int `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Last name of the user
	LastName *string `form:"last_name,omitempty" json:"last_name,omitempty" xml:"last_name,omitempty"`
	// Avatar of user
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
	// Phone number of the user
	Salt *string `form:"salt,omitempty" json:"salt,omitempty" xml:"salt,omitempty"`
}

// GetUserByEmail gets a user by email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	var user User

	err := db.QueryRow("SELECT first_name, last_name, email, password, salt FROM users WHERE email = $1", email).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Salt)

	return &user, err
}

// AddUserToDatabase creates a new user
func AddUserToDatabase(db *sql.DB, firstName, lastName, email, password string) error {
	const sqlstr = `
    INSERT INTO users (
        first_name,
        last_name,
        email,
        password,
        salt
    ) VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) returning id
    `
	salt := crypto.GenerateSalt()
	hashedPassword := crypto.HashPassword(password, salt)
	var id int
	err := db.QueryRow(sqlstr, firstName, lastName, email, hashedPassword, salt).Scan(&id)
	return err
}

func CheckEmailExists(db *sql.DB, email string) (bool, error) {
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}