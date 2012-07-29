package main

import (
  _ "github.com/bmizerany/pq"
  "database/sql"
  "html/template"
  "log"
  "net/http"
)

var templates = template.Must( template.ParseGlob("views/*.html") )

type Bookmark struct {
  Id int
  Name string
  Url string
}

// Open a connection to our database
// return `db`, `*sql.DB` - a pointer to the sql's package DB struct
func dbconnect() *sql.DB {
  db, err := sql.Open("postgres", "user=davidglivar dbname=bookmarked sslmode=disable")
  if err != nil {
    log.Println("Error: Cannot establish connection to database.")
    log.Println(err)
  }
  return db
}

// Close an active connection to our database
// param `db`, `*sql.DB` - a pointer to the sql's package DB struct
func dbclose(db *sql.DB) {
  err := db.Close()
  if err != nil {
    log.Println("Error: Cannot close connection to database.")
    log.Println(err)
  }
}

func getBookmarks() []Bookmark {
  bookmarks := make([]Bookmark, 0)
  db := dbconnect()
  rows, err := db.Query("SELECT * FROM bookmarks")
  if err != nil {
    log.Println("Error Querying database.")
    log.Println(err)
  }

  for rows.Next() {
    var id int
    var name, url string
    err := rows.Scan(&id, &name, &url)
    if err != nil {
      log.Println("Error during row.Scan()")
      log.Println(err)
    }
    bookmark := Bookmark{Id: id, Name: name, Url: url}
    bookmarks = append(bookmarks, bookmark)
  }
  dbclose(db)
  return bookmarks
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  name := r.FormValue("name")
  url := r.FormValue("url")
  db := dbconnect()
  _, err := db.Exec("INSERT INTO bookmarks (name, url) VALUES ($1, $2)", name, url)
  if err != nil {
    log.Println("ERROR! Transaction aborted!")
    log.Println(err)
    return
  }
  dbclose(db)
  // http.Redirect(w, r, "/", http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  tmplData := map[string]interface{} {
    "Bookmarks": getBookmarks(),
  }
  err := templates.ExecuteTemplate(w, "index.html", tmplData)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func main() {
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/create_bookmark", saveHandler)
  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
