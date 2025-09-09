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
//lesson.exe
//
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

	//timeperiods
	mux.HandleFunc("GET /timeperiods/create", timePeriodCreateTemplate)
	mux.HandleFunc("GET /timeperiods/search", timePeriodSearchTemplate)
	mux.HandleFunc("GET /timeperiods/modify/{id}", func(w http.ResponseWriter, r *http.Request) { timePeriodModifyTemplate(w, r, db) })
	mux.HandleFunc("POST /create_timeperiod", func(w http.ResponseWriter, r *http.Request) { createTimePeriod(w, r, db) })
	mux.HandleFunc("POST /search_timeperiod", func(w http.ResponseWriter, r *http.Request) { searchTimePeriod(w, r, db) })
	mux.HandleFunc("POST /modify_timeperiod", func(w http.ResponseWriter, r *http.Request) { modifyTimePeriod(w, r, db) })

	//schools
	mux.HandleFunc("GET /schools/create", schoolCreateTemplate)
	mux.HandleFunc("GET /schools/search", schoolSearchTemplate)
	mux.HandleFunc("GET /schools/modify/{id}", func(w http.ResponseWriter, r *http.Request) { schoolModifyTemplate(w, r, db) })
	mux.HandleFunc("POST /create_school", func(w http.ResponseWriter, r *http.Request) { createSchool(w, r, db) })
	mux.HandleFunc("POST /search_school", func(w http.ResponseWriter, r *http.Request) { searchSchool(w, r, db) })
	mux.HandleFunc("POST /modify_school", func(w http.ResponseWriter, r *http.Request) { modifySchool(w, r, db) })

	//students
	mux.HandleFunc("GET /students/create", func(w http.ResponseWriter, r *http.Request) { studentCreateTemplate(w, r, db) })
	mux.HandleFunc("GET /students/search", studentSearchTemplate)
	mux.HandleFunc("GET /students/modify/{id}", func(w http.ResponseWriter, r *http.Request) { studentModifyTemplate(w, r, db) })
	mux.HandleFunc("POST /create_student", func(w http.ResponseWriter, r *http.Request) { createStudent(w, r, db) })
	mux.HandleFunc("POST /search_student", func(w http.ResponseWriter, r *http.Request) { searchStudent(w, r, db) })
	mux.HandleFunc("POST /modify_student", func(w http.ResponseWriter, r *http.Request) { modifyStudent(w, r, db) })

	//instructors
	mux.HandleFunc("GET /instructors/create", instructorCreateTemplate)
	mux.HandleFunc("GET /instructors/search", instructorSearchTemplate)
	mux.HandleFunc("GET /instructors/modify/{id}", func(w http.ResponseWriter, r *http.Request) { instructorModifyTemplate(w, r, db) })
	mux.HandleFunc("POST /create_instructor", func(w http.ResponseWriter, r *http.Request) { createInstructor(w, r, db) })
	mux.HandleFunc("POST /search_instructor", func(w http.ResponseWriter, r *http.Request) { searchInstructor(w, r, db) })
	mux.HandleFunc("POST /modify_instructor", func(w http.ResponseWriter, r *http.Request) { modifyInstructor(w, r, db) })

	//instructor availability
	mux.HandleFunc("GET /instructoravail/search", func(w http.ResponseWriter, r *http.Request) { instructorAvailSearchTemplate(w, r, db) })
	mux.HandleFunc("POST /search_instructor_avail", func(w http.ResponseWriter, r *http.Request) { searchInstructorAvail(w, r, db) })
	mux.HandleFunc("POST /update_instructor_avail", func(w http.ResponseWriter, r *http.Request) { updateInstructorAvail(w, r, db) })

	fmt.Printf("Server listening to :8080 \n")
	http.ListenAndServe(":8080", mux)
}

func handleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	myMessage := ""
	d := struct {
		TheMessage string
	}{
		TheMessage: myMessage,
	}

	tpl.ExecuteTemplate(w, "index.gohtml", d)
}

func handleError(w http.ResponseWriter, err error, errMsg string) {
	if err != nil {
		d := struct {
			TheMessage string
		}{
			TheMessage: errMsg,
		}

		tpl.ExecuteTemplate(w, "index.gohtml", d)

		fmt.Printf("Error: %s", err.Error())
		return
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Error: %s", err.Error())
		return
	}
}
