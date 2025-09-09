package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func studentCreateTemplate(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	zListBox := getSchoolListBox(db, 0)

	d := struct {
		SchoolList []MyListBox
	}{
		SchoolList: zListBox,
	}

	tpl.ExecuteTemplate(w, "student_create.gohtml", d)
}

func studentSearchTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {

	d := struct {
		YourMessage string
	}{
		YourMessage: "",
	}

	tpl.ExecuteTemplate(w, "student_search.gohtml", d)
}

func studentModifyTemplate(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	yMessage := ""
	pFirstName := ""
	pLastName := ""
	pEmail := ""
	pID := 0
	pSchoolID := 0
	pActive := false

	id, err := strconv.Atoi(r.PathValue("id"))
	studentCheckErr(err)
	students := searchForstudentByID(db, id)

	for index, myStudent := range students {
		if index == 0 {
			pFirstName = myStudent.First_Name
			pLastName = myStudent.Last_Name
			pEmail = myStudent.Email
			pID = myStudent.ID
			pSchoolID = myStudent.School_ID
			pActive = myStudent.Active
		}
	}

	zListBox := getSchoolListBox(db, pSchoolID)

	d := struct {
		YourMessage string
		ID          int
		First_Name  string
		Last_Name   string
		Email       string
		School_ID   int
		Active      bool
		SchoolList  []MyListBox
	}{
		YourMessage: yMessage,
		ID:          pID,
		First_Name:  pFirstName,
		Last_Name:   pLastName,
		Email:       pEmail,
		School_ID:   pSchoolID,
		Active:      pActive,
		SchoolList:  zListBox,
	}

	tpl.ExecuteTemplate(w, "student_modify.gohtml", d)
}

func createStudent(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	pActive := false
	if r.FormValue("form_studentactive") == "TRUE" {
		pActive = true
	}

	schoolList := strings.Split(r.FormValue("form_schools"), "-")
	pSchoolName := schoolList[0]
	pSchoolID, err := strconv.Atoi(schoolList[1])

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	newStudent := Student{
		First_Name: r.FormValue("form_firstname"),
		Last_Name:  r.FormValue("form_lastname"),
		Email:      r.FormValue("form_email"),
		School_ID:  pSchoolID,
		Active:     pActive,
	}

	pID := 0

	cacheMutex.Lock()
	students := insertStudent(db, newStudent, pSchoolName)
	cacheMutex.Unlock()

	for index, myStudent := range students {
		if index == 0 {
			pID = myStudent.ID
		}
	}

	yMessage := fmt.Sprintf("Student %d Created", pID)

	d := struct {
		StudentList []Student
		YourMessage string
		Low         int
		PrevDisplay bool
		NextDisplay bool
		SearchFor   string
	}{
		StudentList: students,
		YourMessage: yMessage,
		Low:         0,
		PrevDisplay: false,
		NextDisplay: false,
		SearchFor:   "",
	}

	tpl.ExecuteTemplate(w, "student_list.gohtml", d)
}

func modifyStudent(
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
		deleteStudent(db, id)
		cacheMutex.Unlock()

		yMessage = fmt.Sprintf("Student %d Deleted", id)
	} else {
		yActive := false
		if r.FormValue("form_studentactive") == "TRUE" {
			yActive = true
		}

		schoolList := strings.Split(r.FormValue("form_schools"), "-")
		ySchoolID, err := strconv.Atoi(schoolList[1])

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		zStudent := Student{
			ID:         id,
			First_Name: r.FormValue("form_firstname"),
			Last_Name:  r.FormValue("form_lastname"),
			Email:      r.FormValue("form_email"),
			School_ID:  ySchoolID,
			Active:     yActive,
		}

		cacheMutex.Lock()
		updateStudent(db, zStudent)
		cacheMutex.Unlock()

		yMessage = fmt.Sprintf("Student %d Modified", id)
	}

	cacheMutex.Lock()
	students := searchForstudentByID(db, id)
	cacheMutex.Unlock()

	d := struct {
		StudentList []Student
		YourMessage string
		Low         int
		PrevDisplay bool
		NextDisplay bool
		SearchFor   string
	}{
		StudentList: students,
		YourMessage: yMessage,
		Low:         0,
		PrevDisplay: false,
		NextDisplay: false,
		SearchFor:   "",
	}

	tpl.ExecuteTemplate(w, "student_list.gohtml", d)
}

func searchStudent(
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

	NewSearchFor := getSearchFor(r.FormValue("form_button"), r.FormValue("form_studentsearch"), r.FormValue("form_prev_next_searchfor"))

	cacheMutex.Lock()
	maxCount = countForstudent(db, NewSearchFor)
	cacheMutex.Unlock()

	zPage := getPrevNext(r.FormValue("form_button"), low, maxCount, rows)

	cacheMutex.Lock()
	students := searchForstudent(db, NewSearchFor, zPage.NewLow, rows)
	cacheMutex.Unlock()

	yMessage := fmt.Sprintf("Searched for %s", NewSearchFor)

	d := struct {
		StudentList []Student
		YourMessage string
		Low         int
		PrevDisplay bool
		NextDisplay bool
		SearchFor   string
	}{
		StudentList: students,
		YourMessage: yMessage,
		Low:         zPage.NewLow,
		PrevDisplay: zPage.PrevDisplay,
		NextDisplay: zPage.NextDisplay,
		SearchFor:   NewSearchFor,
	}

	tpl.ExecuteTemplate(w, "student_list.gohtml", d)
}

func studentCheckErr(err error) {
	if err != nil {
		//d := struct {
		//	YourMessage string
		//}{
		//	YourMessage: "Error",
		//}

		//tpl.ExecuteTemplate(w, "student_search.gohtml", d)

		log.Fatal(err)
		//fmt.Printf("Error: %s", err.Error())
		return
	}
}
