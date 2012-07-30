package main

import (
	"database/sql"
	"errors"
	_ "github.com/bmizerany/pq"
	"html/template"
	"log"
	"net/http"
)

// templates are accessed without the 'views/' prefix
// eg, err := templates.ExecuteTemplate(w, "index.html", tmplData)
var templates = template.Must(template.ParseGlob("views/*.html"))

// Bookmark is our sole model
type Bookmark struct {
	Id   int
	Name string
	Url  string
}

// Open a connection to our database
// return db *sql.DB - a pointer to the sql's package DB struct
func DBconnect() (db *sql.DB) {
	db, err := sql.Open("postgres", "user=davidglivar dbname=bookmarked sslmode=disable")
	if err != nil {
		log.Println("Error: Cannot establish connection to database.")
		log.Println(err)
	}
	return
}

// Close an active connection to our database
// param db *sql.DB - a pointer to the sql's package DB struct
func DBclose(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Error: Cannot close connection to database.")
		log.Println(err)
	}
}

// Bookmarks() - a getter method
// return bookmarks []Bookmark - a slice of Bookmark structs
func Bookmarks() (bookmarks []Bookmark) {
	db := DBconnect()
	defer DBclose(db)
	rows, err := db.Query("SELECT * FROM bookmarks")
	if err != nil {
		log.Println("Error Querying database.")
		log.Println(err)
		return
	}

	bookmarks = make([]Bookmark, 0)
	for rows.Next() {
		var id int
		var name, url string
		err := rows.Scan(&id, &name, &url)
		if err != nil {
			log.Println("Error during row.Scan()")
			log.Println(err)
			return
		}
		bookmark := Bookmark{Id: id, Name: name, Url: url}
		bookmarks = append(bookmarks, bookmark)
	}
	return
}

// Validator() queries the database to ensure we are not
// about to write a duplicate entry.
// return error
func Validator(name, url string) error {
	var id int
	var n, u string
	db := DBconnect()
	defer DBclose(db)
	row := db.QueryRow("SELECT * FROM bookmarks WHERE name = $1 AND url = $2", name, url)
	err := row.Scan(&id, &n, &u)
	if err != nil {
		log.Println(err)
		return nil
	}
	// found a match, throw error
	return errors.New("Match found, close connection.")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	url := r.FormValue("url")

	validationerror := Validator(name, url)
	if validationerror != nil {
		log.Println(validationerror)
		log.Println(w.Header())
		return
	}

	db := DBconnect()
	defer DBclose(db)
	_, err := db.Exec("INSERT INTO bookmarks (name, url) VALUES ($1, $2)", name, url)
	if err != nil {
		log.Println("ERROR! Could not create record.")
		log.Println(err)
		return
	}
	return
}

func destroyHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	db := DBconnect()
	defer DBclose(db)
	_, err := db.Exec("DELETE FROM bookmarks WHERE id = $1", id)
	if err != nil {
		log.Println("ERROR! Could not delete record.")
		log.Println(err)
		return
	}
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmplData := map[string]interface{}{
		"Bookmarks": Bookmarks(),
	}
	err := templates.ExecuteTemplate(w, "index.html", tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create_bookmark", saveHandler)
	http.HandleFunc("/delete_bookmark", destroyHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
