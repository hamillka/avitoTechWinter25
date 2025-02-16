package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	Username string
	Password string
}

func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/avito_shop_service?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	numUsers := 100000
	users := generateUsers(numUsers)

	err = insertUsers(db, users)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Успешно добавлено %d пользователей в базу данных\n", numUsers)
}

func generateUsers(num int) []User {
	var users []User
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < num; i++ {
		username := fmt.Sprintf("user%d", i+1)
		password := fmt.Sprintf("password%d", i+1)
		users = append(users, User{Username: username, Password: password})
	}

	return users
}

func insertUsers(db *sql.DB, users []User) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	stmt, err := tx.Prepare("INSERT INTO users (username, password) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	for _, user := range users {
		_, err := stmt.Exec(user.Username, user.Password)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
