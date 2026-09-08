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
	State      string
	IsLoggedIn bool
	Score      int
}

// saving user to the database
func (u User) Save() error {
	// hash the password
	hashPassword, err := utility.HashPassword(u.Password)
	if err != nil {
		return err
	}

	// Postgres supports RETURNING; MySQL uses LastInsertId
	if database.Currentdb == database.Postgres {
		query := `INSERT INTO users(name, email, password, isloggedin) VALUES($1,$2,$3,$4) RETURNING id`
		query = database.NormalizeQuery(database.Currentdb, query)
		stmt, err := database.Database.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()
		row := stmt.QueryRow(u.Name, u.Email, hashPassword, false)
		if err := row.Scan(&u.ID); err != nil {
			return err
		}
		return nil
	}

	query := `INSERT INTO users(name, email, password, isloggedin) VALUES($1,$2,$3,$4)`
	query = database.NormalizeQuery(database.Currentdb, query)
	stmt, err := database.Database.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(u.Name, u.Email, hashPassword, false)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = lastID
	return nil
}

func (u *User) ValidateUser() error {
	query := "SELECT id, name, password, isloggedin FROM users WHERE email = $1"
	query = database.NormalizeQuery(database.Currentdb, query)
	row := database.Database.QueryRow(query, u.Email)

	var hashedPassword string
	var loggedIn bool

	err := row.Scan(&u.ID, &u.Name, &hashedPassword, &loggedIn)
	if err != nil {
		return errors.New("credentials couldn't be read")
	}

	u.IsLoggedIn = loggedIn

	passWordIsValid := utility.CheckPasswordHash(u.Password, hashedPassword)
	if !passWordIsValid {
		return errors.New("credentials is invalid, password you entered not correct")
	}

	return nil

}

// get all users from database (every registered user)
func GetAllUsers() ([]User, error) {
	query := "SELECT id, name, email, password, isloggedin FROM users"
	rows, err := database.Database.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsLoggedIn)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserByID(id int64) (*User, error) {
	query := "SELECT id, name, email, password, isloggedin, score FROM users WHERE id = $1"
	query = database.NormalizeQuery(database.Currentdb, query)
	row := database.Database.QueryRow(query, id)

	var u User
	var loggedIn bool

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &loggedIn, &u.Score)
	if err != nil {
		return nil, err
	}

	u.IsLoggedIn = loggedIn
	return &u, nil

}

// Add to user's score (can be negative to subtract)
func (u *User) AddToScore(point int) error {
	query := "UPDATE users SET score = score + $1 WHERE id = $2"
	query = database.NormalizeQuery(database.Currentdb, query)

	_, err := database.Database.Exec(query, point, u.ID)
	if err != nil {
		return err
	}

	u.Score += point
	return nil
}

func (u *User) SetLoggedIn(loggedIn bool) error {
	query := "UPDATE users SET isloggedin = $1 WHERE id = $2"
	query = database.NormalizeQuery(database.Currentdb, query)

	_, err := database.Database.Exec(query, loggedIn, u.ID)

	if err != nil {
		return err
	}

	u.IsLoggedIn = loggedIn
	return nil
}

func (u *User) UpdateEmail(newEmail string) error {
	query := "UPDATE users SET email = $1 WHERE id = $2"
	query = database.NormalizeQuery(database.Currentdb, query)
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
	query := "UPDATE users SET password = $1 WHERE id = $2"
	query = database.NormalizeQuery(database.Currentdb, query)
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
