package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func schoolCreateTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {
	tpl.ExecuteTemplate(w, "school_create.gohtml", nil)
}

func schoolSearchTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {

	d := struct {
		YourMessage string
	}{
		YourMessage: "",
	}

	tpl.ExecuteTemplate(w, "school_search.gohtml", d)
}

func schoolModifyTemplate(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	yMessage := ""
	pSchoolName := ""
	pID := 0
	pActive := false

	id, err := strconv.Atoi(r.PathValue("id"))
	schoolCheckErr(err)
	schools := searchForschoolByID(db, id)

	for index, mySchool := range schools {
		if index == 0 {
			pSchoolName = mySchool.School_Name
			pID = mySchool.ID
			pActive = mySchool.Active
		}
	}

	d := struct {
		YourMessage string
		ID          int
		School_Name string
		Active      bool
	}{
		YourMessage: yMessage,
		ID:          pID,
		School_Name: pSchoolName,
		Active:      pActive,
	}

	tpl.ExecuteTemplate(w, "school_modify.gohtml", d)
}

func createSchool(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	pActive := false
	if r.FormValue("form_schoolactive") == "TRUE" {
		pActive = true
	}

	newSchool := School{
		School_Name: r.FormValue("form_schoolname"),
		Active:      pActive,
	}

	pID := 0

	cacheMutex.Lock()
	schools := insertSchool(db, newSchool)
	cacheMutex.Unlock()

	for index, mySchool := range schools {
		if index == 0 {
			pID = mySchool.ID
		}
	}

	yMessage := fmt.Sprintf("School %d Created", pID)

	d := struct {
		SchoolList  []School
		YourMessage string
		Low         int
		PrevDisplay bool
		NextDisplay bool
		SearchFor   string
	}{
		SchoolList:  schools,
		YourMessage: yMessage,
		Low:         0,
		PrevDisplay: false,
		NextDisplay: false,
		SearchFor:   "",
	}

	tpl.ExecuteTemplate(w, "school_list.gohtml", d)
}

func modifySchool(
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
		deleteSchool(db, id)
		cacheMutex.Unlock()

		yMessage = fmt.Sprintf("School %d Deleted", id)
	} else {
		pActive := false
		if r.FormValue("form_schoolactive") == "TRUE" {
			pActive = true
		}

		zSchool := School{
			ID:          id,
			School_Name: r.FormValue("form_schoolname"),
			Active:      pActive,
		}

		cacheMutex.Lock()
		updateSchool(db, zSchool)
		cacheMutex.Unlock()

		yMessage = fmt.Sprintf("School %d Modified", id)
	}

	cacheMutex.Lock()
	schools := searchForschoolByID(db, id)
	cacheMutex.Unlock()

	d := struct {
		SchoolList  []School
		YourMessage string
		Low         int
		PrevDisplay bool
		NextDisplay bool
		SearchFor   string
	}{
		SchoolList:  schools,
		YourMessage: yMessage,
		Low:         0,
		PrevDisplay: false,
		NextDisplay: false,
		SearchFor:   "",
	}

	tpl.ExecuteTemplate(w, "school_list.gohtml", d)
}

func searchSchool(
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

	NewSearchFor := getSearchFor(r.FormValue("form_button"), r.FormValue("form_schoolsearch"), r.FormValue("form_prev_next_searchfor"))

	cacheMutex.Lock()
	maxCount = countForschool(db, NewSearchFor)
	cacheMutex.Unlock()

	zPage := getPrevNext(r.FormValue("form_button"), low, maxCount, rows)

	cacheMutex.Lock()
	schools := searchForschool(db, NewSearchFor, zPage.NewLow, rows)
	cacheMutex.Unlock()

	yMessage := fmt.Sprintf("Searched for %s", NewSearchFor)

	d := struct {
		SchoolList  []School
		YourMessage string
		Low         int
		PrevDisplay bool
		NextDisplay bool
		SearchFor   string
	}{
		SchoolList:  schools,
		YourMessage: yMessage,
		Low:         zPage.NewLow,
		PrevDisplay: zPage.PrevDisplay,
		NextDisplay: zPage.NextDisplay,
		SearchFor:   NewSearchFor,
	}

	tpl.ExecuteTemplate(w, "school_list.gohtml", d)
}

func schoolCheckErr(err error) {
	if err != nil {
		//d := struct {
		//	YourMessage string
		//}{
		//	YourMessage: "Error",
		//}

		//tpl.ExecuteTemplate(w, "school_search.gohtml", d)

		log.Fatal(err)
		//fmt.Printf("Error: %s", err.Error())
		return
	}
}
