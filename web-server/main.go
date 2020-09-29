package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"database/sql"
	"log"
	json "encoding/json"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/bookstore?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Connected")

	// seed := `
	// INSERT INTO books
  // (isbn, title, author, price)
	// VALUES
	// ($1, $2, $3, $4)
	// `
	// id := 0
	// err = db.QueryRow(seed, 12345, "testtitle", "zeger de vos", 40).Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("New record created with id:", id)
}

type Book struct {
	Isbn string
	Title string
	Author string
	Price float32
}

func health(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "hello world")
}

func getBooks(w http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	getBooks := `SELECT * FROM books`
	rows, err := db.Query(getBooks)
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	bks := make([]Book, 0)

	for rows.Next() {
		bk := Book{}
		err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		bks = append(bks, bk)
	}
	if err := rows.Err(); err != nil {
		return
	}
	if err := json.NewEncoder(w).Encode(bks); err != nil {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/", health)
	http.HandleFunc("/books", getBooks)
	http.ListenAndServe(":8080", nil)
}
