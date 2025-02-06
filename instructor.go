package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func instructorCreateTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {
	tpl.ExecuteTemplate(w, "instructor_create.gohtml", nil)
}

func instructorSearchTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {

	d := struct {
		YourMessage string
	}{
		YourMessage: "",
	}

	tpl.ExecuteTemplate(w, "instructor_search.gohtml", d)
}

func instructorModifyTemplate(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	yMessage := ""
	pInstructorName := ""
	pID := 0
	pActive := false

	id, err := strconv.Atoi(r.PathValue("id"))
	instructorCheckErr(err)
	instructors := searchForinstructorByID(db, id)

	for index, myInstructor := range instructors {
		if index == 0 {
			pInstructorName = myInstructor.Instructor_Name
			pID = myInstructor.ID
			pActive = myInstructor.Active
		}
	}

	d := struct {
		YourMessage     string
		ID              int
		Instructor_Name string
		Active          bool
	}{
		YourMessage:     yMessage,
		ID:              pID,
		Instructor_Name: pInstructorName,
		Active:          pActive,
	}

	tpl.ExecuteTemplate(w, "instructor_modify.gohtml", d)
}

func createInstructor(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	pActive := false
	if r.FormValue("form_instructoractive") == "TRUE" {
		pActive = true
	}

	newInstructor := Instructor{
		Instructor_Name: r.FormValue("form_instructorname"),
		Active:          pActive,
	}

	pID := 0

	cacheMutex.Lock()
	instructors := insertInstructor(db, newInstructor)
	cacheMutex.Unlock()

	for index, myInstructor := range instructors {
		if index == 0 {
			pID = myInstructor.ID
		}
	}

	yMessage := fmt.Sprintf("Instructor %d Created", pID)

	d := struct {
		InstructorList []Instructor
		YourMessage    string
		Low            int
		PrevDisplay    bool
		NextDisplay    bool
		SearchFor      string
	}{
		InstructorList: instructors,
		YourMessage:    yMessage,
		Low:            0,
		PrevDisplay:    false,
		NextDisplay:    false,
		SearchFor:      "",
	}

	tpl.ExecuteTemplate(w, "instructor_list.gohtml", d)
}

func modifyInstructor(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {

	fmt.Printf("Form Button = %s \n", r.FormValue("form_button"))

	id, err := strconv.Atoi(r.FormValue("form_id"))

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	yMessage := ""

	if r.FormValue("form_button") == "Delete" {
		cacheMutex.Lock()
		deleteInstructor(db, id)
		cacheMutex.Unlock()

		yMessage = fmt.Sprintf("Instructor %d Deleted", id)
	} else {
		pActive := false
		if r.FormValue("form_instructoractive") == "TRUE" {
			pActive = true
		}

		zInstructor := Instructor{
			ID:              id,
			Instructor_Name: r.FormValue("form_instructorname"),
			Active:          pActive,
		}

		cacheMutex.Lock()
		updateInstructor(db, zInstructor)
		cacheMutex.Unlock()

		yMessage = fmt.Sprintf("Instructor %d Modified", id)
	}

	cacheMutex.Lock()
	instructors := searchForinstructorByID(db, id)
	cacheMutex.Unlock()

	d := struct {
		InstructorList []Instructor
		YourMessage    string
		Low            int
		PrevDisplay    bool
		NextDisplay    bool
		SearchFor      string
	}{
		InstructorList: instructors,
		YourMessage:    yMessage,
		Low:            0,
		PrevDisplay:    false,
		NextDisplay:    false,
		SearchFor:      "",
	}

	tpl.ExecuteTemplate(w, "instructor_list.gohtml", d)
}

func searchInstructor(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	low := 0
	rows := 100
	maxCount := 0

	low, err := strconv.Atoi(r.FormValue("form_low"))

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	NewSearchFor := getSearchFor(r.FormValue("form_button"), r.FormValue("form_instructorsearch"), r.FormValue("form_prev_next_searchfor"))

	cacheMutex.Lock()
	maxCount = countForinstructor(db, NewSearchFor)
	cacheMutex.Unlock()

	zPage := getPrevNext(r.FormValue("form_button"), low, maxCount, rows)

	cacheMutex.Lock()
	instructors := searchForinstructor(db, NewSearchFor, zPage.NewLow, rows)
	cacheMutex.Unlock()

	yMessage := fmt.Sprintf("Searched for %s", NewSearchFor)

	d := struct {
		InstructorList []Instructor
		YourMessage    string
		Low            int
		PrevDisplay    bool
		NextDisplay    bool
		SearchFor      string
	}{
		InstructorList: instructors,
		YourMessage:    yMessage,
		Low:            zPage.NewLow,
		PrevDisplay:    zPage.PrevDisplay,
		NextDisplay:    zPage.NextDisplay,
		SearchFor:      NewSearchFor,
	}

	tpl.ExecuteTemplate(w, "instructor_list.gohtml", d)
}

func instructorCheckErr(err error) {
	if err != nil {
		//d := struct {
		//	YourMessage string
		//}{
		//	YourMessage: "Error",
		//}

		//tpl.ExecuteTemplate(w, "instructor_search.gohtml", d)

		log.Fatal(err)
		//fmt.Printf("Error: %s", err.Error())
		return
	}
}
