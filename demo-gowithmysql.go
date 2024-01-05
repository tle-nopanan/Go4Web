package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func creatingTable(db *sql.DB) {
	query := `CREATE TABLE users (
		id INT AUTO_INCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME,
		PRIMARY KEY (id)
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}

}

func Insert(db *sql.DB) {
	var username string
	var password string
	fmt.Scan(&username)
	fmt.Scan(&password)
	createdAt := time.Now()

	result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	fmt.Println(id)

}

func query(db *sql.DB) {
	var (
		id         int
		coursename string
		price      float64
		instructor string
	)
	for {
		var inputID int
		fmt.Println("Input id (or type 'exit' to quit)")
		_, err := fmt.Scan(&inputID)
		if err != nil {
			log.Fatal(err)
		}

		if inputID == -1 { // Set the exit condition (you can change this to any condition)
			fmt.Println("Exiting loop...")
			break
		}

		query := "SELECT id, coursename, price, instructor FROM onlinecourse WHERE id = ?"
		if err := db.QueryRow(query, inputID).Scan(&id, &coursename, &price, &instructor); err != nil {
			if err == sql.ErrNoRows {
				log.Println("ID not found in database. Please input another ID.")
				continue // Continue to next iteration to ask for another ID
			}
			log.Fatal(err)
		}
		fmt.Println(id, coursename, price, instructor)
	}
}

func delete(db *sql.DB) {
	var deleteid int
	fmt.Scan(&deleteid)
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, deleteid)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		fmt.Println("fail to connect")
	} else {
		fmt.Println("connect success")
	}
	// creatingTable(db)
	// Insert(db)
	delete(db)

}
