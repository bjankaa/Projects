package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var Database *sql.DB

func getEnvOrFile(key string) string {

	fileKey := key + "_FILE"
	if path := os.Getenv(fileKey); path != "" {
		b, err := os.ReadFile(path)
		if err == nil {
			return string(b)
		}
		log.Printf("Failed reading %s: %v; falling back to %s", fileKey, err, key)
	}
	return os.Getenv(key)
}

func InitDatabase() {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "postgres"
	}
	dsn := getEnvOrFile("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL not set")
	}

	if driver == "mysql" {
		Currentdb = MySQL
	} else {
		Currentdb = Postgres
	}

	var datab *sql.DB
	var err error
	maxAttempts := 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		datab, err = sql.Open(driver, dsn)
		if err != nil {
			log.Printf("Can't reach database: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		// Ping with timeout for ready
		pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pingErr := datab.PingContext(pingCtx)
		cancel()
		if pingErr == nil {
			break
		}
		log.Printf("Can't reach database: %v", pingErr)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic("Database could't connect! " + err.Error())
	}

	Database = datab

	// Connection pool configuration
	// setting the maximum number of database connections
	if maxOpenStr := os.Getenv("DB_MAX_OPEN"); maxOpenStr != "" {
		if maxOpen, err := strconv.Atoi(maxOpenStr); err == nil {
			Database.SetMaxOpenConns(maxOpen)
		} else {
			log.Printf("Invalid DB_MAX_OPEN: %v", err)
		}
	}
	// maximum number of unused/idle connections
	if maxIdleStr := os.Getenv("DB_MAX_IDLE"); maxIdleStr != "" {
		if maxIdle, err := strconv.Atoi(maxIdleStr); err == nil {
			Database.SetMaxIdleConns(maxIdle)
		} else {
			log.Printf("Invalid DB_MAX_IDLE: %v", err)
		}
	}

	// maximum lifetime of a connection
	if lifetimeStr := os.Getenv("DB_CONN_MAX_LIFETIME"); lifetimeStr != "" {
		if dur, err := time.ParseDuration(lifetimeStr); err == nil {
			Database.SetConnMaxLifetime(dur)
		} else {
			log.Printf("Invalid DB_CONN_MAX_LIFETIME: %v", err)
		}
	}

	// Create tables for chosen driver (no migrations path)
	if driver == "mysql" {
		err = createTablesMySQL()
	} else {
		err = creatTablesPostgers()
	}

	if err != nil {
		log.Println("Database failure!" + err.Error())
	}

	if err := resetLoginState(); err != nil {
		log.Println("Could not reset login state: " + err.Error())
	}

	fmt.Println("Tables have been created")
}

func resetLoginState() error {
	_, err := Database.Exec(`UPDATE users SET isloggedin = FALSE`)
	return err
}

func creatTablesPostgers() error {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		isloggedin BOOLEAN NOT NULL DEFAULT FALSE,
		score INTEGER NOT NULL DEFAULT 0
	)
	`

	createFinishedGames := `
	CREATE TABLE IF NOT EXISTS finished_games (
		id SERIAL PRIMARY KEY,
		player1_id INTEGER NOT NULL,
		player1_name TEXT NOT NULL,
		player2_id INTEGER NOT NULL,
		player2_name TEXT NOT NULL,
		player3_id INTEGER NOT NULL,
		player3_name TEXT NOT NULL,
		declarer_id INTEGER NOT NULL,
		declarer_name TEXT NOT NULL,
		declarer_win BOOLEAN NOT NULL,
		declarer_points INTEGER NOT NULL,
		defenders_points INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	_, err := Database.Exec(createUserTable)

	if err != nil {
		panic("Database could not connect: user" + err.Error())
	}
	_, err = Database.Exec(createFinishedGames)
	if err != nil {
		panic("Database could not connect: finished_games" + err.Error())
	}
	return err
}

// MySQL variants (placeholder differences: use '?' in queries)
func createTablesMySQL() error {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		isloggedin TINYINT(1) NOT NULL DEFAULT 0,
		score INT NOT NULL DEFAULT 0
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	createFinishedGames := `
	CREATE TABLE IF NOT EXISTS finished_games (
		id INT AUTO_INCREMENT PRIMARY KEY,
		player1_id INT NOT NULL,
		player1_name VARCHAR(255) NOT NULL,
		player2_id INT NOT NULL,
		player2_name VARCHAR(255) NOT NULL,
		player3_id INT NOT NULL,
		player3_name VARCHAR(255) NOT NULL,
		declarer_id INT NOT NULL,
		declarer_name VARCHAR(255) NOT NULL,
		declarer_win TINYINT(1) NOT NULL,
		declarer_points INT NOT NULL,
		defenders_points INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	_, err := Database.Exec(createUserTable)
	if err != nil {
		panic("Database could not connect: user (mysql)" + err.Error())
	}
	_, err = Database.Exec(createFinishedGames)
	if err != nil {
		panic("Database could not connect: finished_games (mysql)" + err.Error())
	}
	return err
}
