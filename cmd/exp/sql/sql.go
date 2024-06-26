package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)
}

func main() {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
		SSLMode:  "disable",
	}

	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			amount INT,
			description TEXT
		);
	`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created!")

	/* Insert user
	name, email := "Bob Calhoun", "bob@calhoun.io"
	row := db.QueryRow(`
		INSERT INTO users(name, email)
		VALUES($1, $2) RETURNING id;
	`, name, email)
	var id int
	err = row.Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User created with id %d\n", id)
	*/

	id := 1
	row := db.QueryRow(`
		SELECT name, email
		FROM users
		WHERE id = $1;
	`, id)
	var name, email string
	err = row.Scan(&name, &email)
	if err == sql.ErrNoRows {
		fmt.Println("Error, no rows!")
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("User %d - Name:%q Email:%q\n", id, name, email)

	/* Insert orders
	userID := 1
	for i := 1; i <= 5; i++ {
		amount := i * 100
		desc := fmt.Sprintf("Fake order #%d", i)
		_, err := db.Exec(`
			INSERT INTO orders(user_id, amount, description)
			VALUES($1, $2, $3)
		`, userID, amount, desc)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Created orders.")
	*/

	type Order struct {
		ID          int
		UserID      int
		Amount      string
		Description string
	}
	var orders []Order
	userID := 1
	rows, err := db.Query(`
		SELECT id, amount, description
		FROM orders
		WHERE user_id = $1;
	`, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		order.UserID = userID
		err := rows.Scan(&order.ID, &order.Amount, &order.Description)
		if err != nil {
			panic(err)
		}
		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Orders: %+v\n", orders)
}
