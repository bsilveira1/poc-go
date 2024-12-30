package main

import (
	"database/sql"
	"log"
	"math/rand"

	// _ "github.com/lib/pq" // postgres driver
	_ "github.com/go-sql-driver/mysql" // mysql driver
)

func main() {
	// PostgreSQL connection (commented)
	// db, err := sql.Open("postgres", "host=postgres user=testuser password=testpassword dbname=testdb sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// MySQL connection (uncommented)
	db, err := sql.Open("mysql", "testuser:testpassword@tcp(localhost:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 100; i++ {
		_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", randomName(), randomEmail())
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Fake data inserted successfully!")
}

func randomName() string {
	names := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	return names[rand.Intn(len(names))]
}

func randomEmail() string {
	domains := []string{"example.com", "test.com", "mail.com"}
	return randomName() + "@" + domains[rand.Intn(len(domains))]
}
