package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"database/sql"
	"log"
	json "encoding/json"
	"strconv"
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
	Isbn int
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

func postBook(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, http.StatusText(400), 400)
	}
	book := Book{}
	book.Isbn, _ = strconv.Atoi(req.FormValue("Isbn"))
	book.Title = req.FormValue("Title")
	book.Author = req.FormValue("Willem")
	f64, err := strconv.ParseFloat(req.FormValue("Price"), 32)
	book.Price = float32(f64)

	// some validation
	if book.Title == "" || book.Author == "" {
		http.Error(res, "Incorrect value for Title or Author", 400)
	}
	_, err = db.Exec(`INSERT INTO books (isbn, title, author, price) VALUES ($1, $2, $3, $4)`, book.Isbn, book.Title, book.Author, book.Price)
	if err != nil {
		log.Fatal(err)
		http.Error(res, http.StatusText(500), 500)
		return
	}
}

func updatePrice(res http.ResponseWriter, req *http.Request) {
	if req.Method != "PUT" {
		http.Error(res, http.StatusText(400), 400)
	}
	keys, ok := req.URL.Query()["isbn"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'isbn' is missing")
		return
	}
	key := keys[0]
	
	book := Book{}
	book.Isbn, _ = strconv.Atoi(key)
	fmt.Println(book.Isbn)
	f64, _ := strconv.ParseFloat(req.FormValue("Price"), 32)
	book.Price = float32(f64)
	updatePrice := `UPDATE books SET PRICE = $1 WHERE Isbn = $2`
	_, err := db.Exec(updatePrice, book.Price, book.Isbn)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", health)
	http.HandleFunc("/book", postBook)
	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/book/price", updatePrice)
	http.ListenAndServe(":8080", nil)
}
