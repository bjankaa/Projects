package models

import (
	"errors"

	"exmaple.com/ulti-restapi/database"
	"exmaple.com/ulti-restapi/utility"
)

type User struct {
	ID         int64
	Name       string
	Email      string `binding:"required"`
	Password   string `binding:"required"`
	State      string `binding:"required"`
	IsLoggedIn bool
}

// saving user to the database
func (u User) Save() error {
	query := `INSERT INTO users(name, email, password, state, isloggedin) VALUES(?,?,?,?,?)`
	stmt, err := database.Database.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	hashPassword, err := utility.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Name, u.Email, hashPassword, u.State, false)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	u.ID = id
	return err
}

func (u *User) ValidateUser() error {
	query := "SELECT id, name, password, state, isloggedin FROM users WHERE email = ?"
	row := database.Database.QueryRow(query, u.Email)

	var hashedPassword string
	var loggedInInt int

	err := row.Scan(&u.ID, &u.Name, &hashedPassword, &u.State, &loggedInInt)
	if err != nil {
		return errors.New("credentials couldn't be read")
	}

	u.IsLoggedIn = loggedInInt == 1

	passWordIsValid := utility.CheckPasswordHash(u.Password, hashedPassword)
	if !passWordIsValid {
		return errors.New("credentials is invalid, password you entered not correct")
	}

	return nil

}

// get all users from database (every registered user)
func GetAllUsers() ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := database.Database.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.State, &user.IsLoggedIn)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserByID(id int64) (*User, error) {
	query := "SELECT id, name, email, password, state, isloggedin FROM users WHERE id = ?"
	row := database.Database.QueryRow(query, id)

	var u User
	var loggedInInt int

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.State, &loggedInInt)
	if err != nil {
		return nil, err
	}

	u.IsLoggedIn = loggedInInt == 1
	return &u, nil

}

func (u *User) SetLoggedIn(loggedIn bool) error {
	query := "UPDATE users SET isloggedin = ? WHERE id = ?"

	_, err := database.Database.Exec(query, loggedIn, u.ID)

	if err != nil {
		return err
	}

	u.IsLoggedIn = loggedIn
	return nil
}

func (u *User) UpdateEmail(newEmail string) error {
	query := "UPDATE users SET email = ? WHERE id = ?"
	stmt, err := database.Database.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newEmail, u.ID)
	if err != nil {
		return err
	}

	u.Email = newEmail
	return nil
}

func (u *User) UpdatePassword(hashedPassword string) error {
	query := "UPDATE users SET password = ? WHERE id = ?"
	stmt, err := database.Database.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hashedPassword, u.ID)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	return nil
}
