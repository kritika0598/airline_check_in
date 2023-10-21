package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "secret"
	hostname = "127.0.0.1:3306"
	dbname   = "test"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

type UserDetails struct {
	name string
	id   int
}

func main() {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		// log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		fmt.Printf("Errors %s pinging DB", err)
		return
	}
	fmt.Printf("Connected to DB %s successfully\n", dbname)

	rows, err := db.QueryContext(ctx, "SELECT name, id FROM users")
	if err != nil {
		log.Fatalf("Error getting rows for query: %s", err)
	}
	defer rows.Close()

	var users []*UserDetails
	for rows.Next() {
		var user UserDetails
		if err := rows.Scan(&user.name, &user.id); err != nil {
			log.Fatalf("Error scanning row: %s", err)
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("Error scanning row: %s", err)
	}

	for _, user := range users {
		fmt.Printf("User Name: %s, ID: %d \n", user.name, user.id)
	}
}
