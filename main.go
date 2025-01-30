package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

//go mod init lesson
//go mod tidy
//SET CGO_ENABLED=1
//go build .
//go run main.go

var tpl *template.Template
var cacheMutex sync.RWMutex

func init() {
	fmt.Printf("init function\n")
	tpl = template.Must(template.ParseGlob("templates/*gohtml"))
}

func main() {
	// Connect to database
	db, err := sql.Open("sqlite3", "./lesson.db")
	checkErr(err)

	// defer close
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("GET /timeperiods", timePeriodTemplate)

	//students
	mux.HandleFunc("GET /students/create", studentCreateTemplate)
	mux.HandleFunc("GET /students/search", studentSearchTemplate)
	mux.HandleFunc("GET /students/modify/{id}", func(w http.ResponseWriter, r *http.Request) { studentModifyTemplate(w, r, db) })
	mux.HandleFunc("POST /create_student", func(w http.ResponseWriter, r *http.Request) { createStudent(w, r, db) })
	mux.HandleFunc("POST /search_student", func(w http.ResponseWriter, r *http.Request) { searchStudent(w, r, db) })
	mux.HandleFunc("POST /modify_student", func(w http.ResponseWriter, r *http.Request) { modifyStudent(w, r, db) })

	//schools
	mux.HandleFunc("GET /schools/create", schoolCreateTemplate)
	mux.HandleFunc("GET /schools/search", schoolSearchTemplate)
	mux.HandleFunc("GET /schools/modify/{id}", func(w http.ResponseWriter, r *http.Request) { schoolModifyTemplate(w, r, db) })
	mux.HandleFunc("POST /create_school", func(w http.ResponseWriter, r *http.Request) { createSchool(w, r, db) })
	mux.HandleFunc("POST /search_school", func(w http.ResponseWriter, r *http.Request) { searchSchool(w, r, db) })
	mux.HandleFunc("POST /modify_school", func(w http.ResponseWriter, r *http.Request) { modifySchool(w, r, db) })

	fmt.Printf("Server listening to :8080 \n")
	http.ListenAndServe(":8080", mux)
}

func handleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func timePeriodTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {
	tpl.ExecuteTemplate(w, "timeperiods.gohtml", nil)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Error: %s", err.Error())
		return
	}
}
